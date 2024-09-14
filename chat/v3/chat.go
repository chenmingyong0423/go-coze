package chat

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/chenmingyong0423/go-coze/common/response"
	jsoniter "github.com/json-iterator/go"

	"github.com/chenmingyong0423/go-coze/common/request"
)

const (
	InternationalChatUrl        = "https://api.coze.com/v3/chat"
	InternationalRetrieveUrl    = "https://api.coze.com/v3/chat/retrieve"
	InternationalMessageListUrl = "https://api.coze.com/v3/chat/message/list"
	InternationalCancelUrl      = "https://api.coze.com/v3/chat/cancel"

	chatUrl               = "https://api.coze.cn/v3/chat"
	retrieveUrl           = "https://api.coze.cn/v3/chat/retrieve"
	messageListUrl        = "https://api.coze.cn/v3/chat/message/list"
	cancelUrl             = "https://api.coze.cn/v3/chat/cancel"
	HeaderAuthorization   = "authorization"
	HeaderContentType     = "Content-Type"
	HeaderApplicationJson = "application/json"
)

type Chat struct {
	botID         string
	userId        string
	authorization string
}

func NewChat(authorization, userID, botID string) *Chat {
	return &Chat{
		botID:         botID,
		userId:        userID,
		authorization: authorization,
	}
}

func (c *Chat) ChatRequest() *CreateRequest {
	return &CreateRequest{
		BotID:  c.botID,
		UserId: c.userId,
		chat:   c,
		Stream: false,
	}
}

func (c *Chat) RetrieveRequest(conversationId string) *RetrieveRequest {
	return &RetrieveRequest{
		chat:           c,
		conversationId: conversationId,
	}
}

func (c *Chat) MessageListRequest(conversationId string) *MessageListRequest {
	return &MessageListRequest{
		chat:           c,
		conversationId: conversationId,
	}
}

func (c *Chat) CancelRequest(conversationId string) *CancelRequest {
	return &CancelRequest{
		chat:           c,
		conversationId: conversationId,
	}
}

type CreateRequest struct {
	chat    *Chat
	timeout time.Duration

	// Optional: Indicate which conversation the dialog is taking place in.
	// 可选的：标识对话发生在哪一次会话中，使用方自行维护此字段。
	conversationId string

	BotID  string `json:"bot_id"`
	UserId string `json:"user_id"`

	// Whether to stream the response to the client.
	// 使用启用流式返回。
	Stream bool `json:"stream"`

	// Optional: The chat history to pass as the context, sorted in ascending order of time.
	// 可选的：作为上下文传递的聊天历史记录，整个列表按时间升序排序。
	AdditionalMessages []request.EnterMessage `json:"additional_messages,omitempty"`

	// Optional: The customized variable in a key-value pair.
	// 可选的：Bot 中定义的变量。在 Bot prompt中设置好变量{{key}}后，可以通过改参数传入变量值。
	CustomVariables map[string]any `json:"custom_variables,omitempty"`

	// Whether to automatically save the history of conversation records.
	// 是否保存本次对话记录。
	AutoSaveHistory bool `json:"auto_save_history,omitempty"`

	// Additional information, typically used to encapsulate some business-related fields. When viewing the details of chat messages, the system will pass through this additional information.
	// 附加信息，通常用于封装一些业务相关的字段。查看对话消息详情时，系统会透传此附加信息。
	MetaData map[string]any `json:"meta_data,omitempty"`

	// 附加参数，通常用于特殊场景下指定一些必要参数供模型判断，例如指定经纬度，并询问 Bot 此位置的天气。
	ExtraParams []string `json:"extra_params,omitempty"`
}

func (r *CreateRequest) WithTimeout(timeout time.Duration) *CreateRequest {
	r.timeout = timeout
	return r
}

func (r *CreateRequest) WithConversationId(conversationId string) *CreateRequest {
	r.conversationId = conversationId
	return r
}

func (r *CreateRequest) AddMessages(additionalMessages ...request.EnterMessage) *CreateRequest {
	r.AdditionalMessages = append(r.AdditionalMessages, additionalMessages...)
	return r
}

func (r *CreateRequest) WithCustomVariables(customVariables map[string]any) *CreateRequest {
	r.CustomVariables = customVariables
	return r
}

func (r *CreateRequest) WithAutoSaveHistory(autoSaveHistory bool) *CreateRequest {
	r.AutoSaveHistory = autoSaveHistory
	return r
}

func (r *CreateRequest) WithMetaData(metaData map[string]any) *CreateRequest {
	r.MetaData = metaData
	return r
}

func (r *CreateRequest) WithExtraParams(extraParams []string) *CreateRequest {
	r.ExtraParams = extraParams
	return r
}

func (r *CreateRequest) Do(ctx context.Context) (*response.DataResponse[*response.Chat], error) {
	r.Stream = false

	body, err := jsoniter.Marshal(r)
	if err != nil {
		return nil, err
	}

	resp := new(response.DataResponse[*response.Chat])

	// 构建查询参数
	params := url.Values{}
	if r.conversationId != "" {
		params.Add("conversation_id", r.conversationId)
	}
	u, err := url.Parse(chatUrl)
	if err != nil {
		return nil, err
	}
	u.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add(HeaderContentType, HeaderApplicationJson)
	req.Header.Add(HeaderAuthorization, fmt.Sprintf("Bearer %s", r.chat.authorization))

	client := http.DefaultClient
	if r.timeout != 0 {
		client.Timeout = r.timeout
	}

	httpResp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	data, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return resp, err
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, &response.HttpErrorResponse{
			Status:     httpResp.Status,
			StatusCode: httpResp.StatusCode,
			Body:       data,
		}
	}
	if err = jsoniter.Unmarshal(data, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}

func (r *CreateRequest) DoStream(ctx context.Context) (<-chan *StreamingResponse, <-chan error) {
	r.Stream = true
	respChan := make(chan *StreamingResponse)
	errChan := make(chan error)

	go func() {
		defer close(respChan)
		defer close(errChan)

		body, err := jsoniter.Marshal(r)
		if err != nil {
			errChan <- err
			return
		}

		params := url.Values{}
		if r.conversationId != "" {
			params.Add("conversation_id", r.conversationId)
		}
		u, err := url.Parse(chatUrl)
		if err != nil {
			errChan <- err
			return
		}
		u.RawQuery = params.Encode()

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(body))
		if err != nil {
			errChan <- err
			return
		}
		req.Header.Add(HeaderContentType, HeaderApplicationJson)
		req.Header.Add(HeaderAuthorization, fmt.Sprintf("Bearer %s", r.chat.authorization))

		client := http.DefaultClient
		if r.timeout != 0 {
			client.Timeout = r.timeout
		}

		httpResp, err := client.Do(req)
		if err != nil {
			errChan <- err
			return
		}
		defer httpResp.Body.Close()

		if httpResp.StatusCode != http.StatusOK {
			data, _ := io.ReadAll(httpResp.Body)
			errChan <- &response.HttpErrorResponse{
				Status:     httpResp.Status,
				StatusCode: httpResp.StatusCode,
				Body:       data,
			}
			return
		}

		scanner := bufio.NewScanner(httpResp.Body)

		sr := &StreamingResponse{}
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "event:") {
				event := strings.TrimSpace(line[len("event:"):])
				if event == "[DONE]" {
					break
				} else {
					sr.Event = event
					continue
				}
			} else if strings.HasPrefix(line, "data:") {
				data := strings.TrimSpace(line[len("data:"):])
				if data == "conversation.chat.failed" {
					var errResp response.BaseResponse
					if err = jsoniter.UnmarshalFromString(data, &errResp); err != nil {
						errChan <- err
						return
					}
					sr.BaseResponse = errResp
					newSr := *sr
					resetStreamResponse(sr)
					respChan <- &newSr
					break
				} else {
					if strings.Contains(sr.Event, "chat") {
						var chatResp response.Chat
						if err = jsoniter.UnmarshalFromString(data, &chatResp); err != nil {
							errChan <- err
							return
						}
						sr.Chat = &chatResp
						newSr := *sr
						resetStreamResponse(sr)
						respChan <- &newSr
					} else if strings.Contains(sr.Event, "message") {
						var messageResp response.Message
						if err = jsoniter.UnmarshalFromString(data, &messageResp); err != nil {
							errChan <- err
							return
						}
						sr.Message = &messageResp
						newSr := *sr
						resetStreamResponse(sr)
						respChan <- &newSr
					}
				}
			} else if line != "" {
				var resp response.BaseResponse
				if err = jsoniter.UnmarshalFromString(line, &resp); err != nil {
					errChan <- err
					return
				}
				sr.BaseResponse = resp
				newSr := *sr
				resetStreamResponse(sr)
				respChan <- &newSr
				return
			}
		}

		if err = scanner.Err(); err != nil {
			errChan <- err
			return
		}
	}()

	return respChan, errChan
}

// Reset 如果你想复用该对象，建议调用该方法重置。
func (r *CreateRequest) Reset() {
	r.Stream = false
	r.AdditionalMessages = nil
	r.CustomVariables = nil
	r.AutoSaveHistory = false
	r.MetaData = nil
	r.ExtraParams = nil
}

func resetStreamResponse(sr *StreamingResponse) {
	sr.Event = ""
	sr.Chat = nil
	sr.Message = nil
	sr.BaseResponse.Code = 0
	sr.BaseResponse.Msg = ""
}

type RetrieveRequest struct {
	chat    *Chat
	timeout time.Duration

	conversationId string
}

func (r *RetrieveRequest) WithTimeout(timeout time.Duration) *RetrieveRequest {
	r.timeout = timeout
	return r
}

func (r *RetrieveRequest) Do(ctx context.Context, chatId string) (*response.DataResponse[*response.Chat], error) {
	resp := new(response.DataResponse[*response.Chat])

	params := url.Values{}
	params.Add("conversation_id", r.conversationId)
	params.Add("chat_id", chatId)

	u, err := url.Parse(retrieveUrl)
	if err != nil {
		return nil, err
	}
	u.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add(HeaderContentType, HeaderApplicationJson)
	req.Header.Add(HeaderAuthorization, fmt.Sprintf("Bearer %s", r.chat.authorization))

	client := http.DefaultClient
	if r.timeout != 0 {
		client.Timeout = r.timeout
	}

	httpResp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	data, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return resp, err
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, &response.HttpErrorResponse{
			Status:     httpResp.Status,
			StatusCode: httpResp.StatusCode,
			Body:       data,
		}
	}
	if err = jsoniter.Unmarshal(data, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}

type MessageListRequest struct {
	chat    *Chat
	timeout time.Duration

	conversationId string
}

func (r *MessageListRequest) WithTimeout(timeout time.Duration) *MessageListRequest {
	r.timeout = timeout
	return r
}

func (r *MessageListRequest) Do(ctx context.Context, chatId string) (*response.DataResponse[[]response.Message], error) {
	resp := new(response.DataResponse[[]response.Message])
	resp.Data = make([]response.Message, 0)

	// 构建查询参数
	params := url.Values{}
	params.Add("conversation_id", r.conversationId)
	params.Add("chat_id", chatId)

	u, err := url.Parse(messageListUrl)
	if err != nil {
		return nil, err
	}
	u.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add(HeaderContentType, HeaderApplicationJson)
	req.Header.Add(HeaderAuthorization, fmt.Sprintf("Bearer %s", r.chat.authorization))

	client := http.DefaultClient
	if r.timeout != 0 {
		client.Timeout = r.timeout
	}

	httpResp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	data, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return resp, err
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, &response.HttpErrorResponse{
			Status:     httpResp.Status,
			StatusCode: httpResp.StatusCode,
			Body:       data,
		}
	}
	if err = jsoniter.Unmarshal(data, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}

type CancelRequest struct {
	chat    *Chat
	timeout time.Duration

	conversationId string
}

func (r *CancelRequest) WithTimeout(timeout time.Duration) *CancelRequest {
	r.timeout = timeout
	return r
}

func (r *CancelRequest) Do(ctx context.Context, chatId string) (*response.DataResponse[*response.Chat], error) {
	resp := new(response.DataResponse[*response.Chat])

	params := url.Values{}
	params.Add("conversation_id", r.conversationId)
	params.Add("chat_id", chatId)

	u, err := url.Parse(cancelUrl)
	if err != nil {
		return nil, err
	}
	u.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add(HeaderContentType, HeaderApplicationJson)
	req.Header.Add(HeaderAuthorization, fmt.Sprintf("Bearer %s", r.chat.authorization))

	client := http.DefaultClient
	if r.timeout != 0 {
		client.Timeout = r.timeout
	}

	httpResp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	data, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return resp, err
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, &response.HttpErrorResponse{
			Status:     httpResp.Status,
			StatusCode: httpResp.StatusCode,
			Body:       data,
		}
	}
	if err = jsoniter.Unmarshal(data, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}
