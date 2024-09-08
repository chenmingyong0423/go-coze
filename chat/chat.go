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

	"github.com/chenmingyong0423/go-coze/common/request"
	jsoniter "github.com/json-iterator/go"
)

const (
	InternationalChatUrl = "https://api.coze.com/v3/chat"

	chatUrl               = "https://api.coze.cn/v3/chat"
	HeaderAuthorization   = "Authorization"
	HeaderContentType     = "Content-Type"
	HeaderApplicationJson = "application/json"
)

type Chat struct {
	BotID         string        `json:"bot_id"`
	UserId        string        `json:"user_id"`
	Authorization string        `json:"-"`
	Timeout       time.Duration `json:"-"`

	// Optional: Indicate which conversation the dialog is taking place in.
	// 可选的：标识对话发生在哪一次会话中，使用方自行维护此字段。
	ConversationId string `json:"-"`

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

func NewChat(authorization, userID, botID string) *Chat {
	return &Chat{
		BotID:         botID,
		UserId:        userID,
		Authorization: authorization,
	}
}

func (c *Chat) Chat(ctx context.Context) (*DataResponse, error) {
	if c.Stream {
		return nil, fmt.Errorf("stream request not supported")
	}

	body, err := jsoniter.Marshal(c)
	if err != nil {
		return nil, err
	}

	resp := new(DataResponse)

	// 构建查询参数
	params := url.Values{}
	if c.ConversationId != "" {
		params.Add("conversation_id", c.ConversationId)
	}
	u, err := url.Parse(chatUrl)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add(HeaderContentType, HeaderApplicationJson)
	req.Header.Add(HeaderAuthorization, fmt.Sprintf("Bearer %s", c.Authorization))

	client := http.DefaultClient
	if c.Timeout != 0 {
		client.Timeout = c.Timeout
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
		return nil, &HttpErrorResponse{
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

func (c *Chat) StreamChat(ctx context.Context) (<-chan *StreamingResponse, <-chan error) {
	// 创建返回的通道
	respChan := make(chan *StreamingResponse)
	errChan := make(chan error)

	go func() {
		defer close(respChan)
		defer close(errChan)

		body, err := jsoniter.Marshal(c)
		if err != nil {
			errChan <- err
			return
		}

		// 构建查询参数
		params := url.Values{}
		if c.ConversationId != "" {
			params.Add("conversation_id", c.ConversationId)
		}
		u, err := url.Parse(chatUrl)
		if err != nil {
			errChan <- err
			return
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(body))
		if err != nil {
			errChan <- err
			return
		}
		req.Header.Add(HeaderContentType, HeaderApplicationJson)
		req.Header.Add(HeaderAuthorization, fmt.Sprintf("Bearer %s", c.Authorization))

		client := http.DefaultClient
		if c.Timeout != 0 {
			client.Timeout = c.Timeout
		}

		httpResp, err := client.Do(req)
		if err != nil {
			errChan <- err
			return
		}
		defer httpResp.Body.Close()

		if httpResp.StatusCode != http.StatusOK {
			data, _ := io.ReadAll(httpResp.Body)
			errChan <- &HttpErrorResponse{
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
					var errResp BaseResponse
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
				var resp BaseResponse
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

func resetStreamResponse(sr *StreamingResponse) {
	sr.Event = ""
	sr.Chat = nil
	sr.Message = nil
	sr.BaseResponse.Code = 0
	sr.BaseResponse.Msg = ""
}

func (c *Chat) WithStream(stream bool) *Chat {
	c.Stream = stream
	return c
}

func (c *Chat) WithTimeout(timeout time.Duration) *Chat {
	c.Timeout = timeout
	return c
}

func (c *Chat) WithConversationId(conversationId string) *Chat {
	c.ConversationId = conversationId
	return c
}

func (c *Chat) AddMessages(additionalMessages ...request.EnterMessage) *Chat {
	c.AdditionalMessages = append(c.AdditionalMessages, additionalMessages...)
	return c
}

func (c *Chat) WithCustomVariables(customVariables map[string]any) *Chat {
	c.CustomVariables = customVariables
	return c
}

func (c *Chat) WithAutoSaveHistory(autoSaveHistory bool) *Chat {
	c.AutoSaveHistory = autoSaveHistory
	return c
}

func (c *Chat) WithMetaData(metaData map[string]any) *Chat {
	c.MetaData = metaData
	return c
}

func (c *Chat) WithExtraParams(extraParams []string) *Chat {
	c.ExtraParams = extraParams
	return c
}
