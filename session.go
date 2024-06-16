package coze

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

const (
	InternationalUrl = "https://api.coze.com/open_api/v2/chat"

	url                   = "https://api.coze.cn/open_api/v2/chat"
	HeaderAuthorization   = "Authorization"
	HeaderContentType     = "Content-Type"
	HeaderConnection      = "Connection"
	HeaderAccept          = "Accept"
	HeaderKeepAlive       = "keep-alive"
	HeaderAcceptAll       = "*/*"
	HeaderApplicationJson = "application/json"
)

type Session struct {
	BotID               string        `json:"bot_id"`
	User                string        `json:"user"`
	PersonalAccessToken string        `json:"-"`
	Timeout             time.Duration `json:"-"`

	// Optional: Indicate which conversation the dialog is taking place in.
	// 可选的：标识对话发生在哪一次会话中，使用方自行维护此字段。
	ConversationId string `json:"conversation_id,omitempty"`

	// Whether to stream the response to the client.
	// 使用启用流式返回。
	Stream bool `json:"stream"`

	// The query sent to the bot.
	// 发送给 Bot 的消息内容。
	Query string `json:"query"`

	// Optional: The chat history to pass as the context, sorted in ascending order of time.
	// 可选的：作为上下文传递的聊天历史记录，整个列表按时间升序排序。
	ChatHistory []Message `json:"chat_history,omitempty"`

	// Optional: The customized variable in a key-value pair.
	// 可选的：Bot 中定义的变量。在 Bot prompt中设置好变量{{key}}后，可以通过改参数传入变量值。
	CustomVariables map[string]any `json:"custom_variables,omitempty"`
}

func NewSession(chat *Chat, stream bool) *Session {
	return &Session{
		BotID:               chat.BotID,
		User:                chat.User,
		PersonalAccessToken: chat.PersonalAccessToken,
		Timeout:             chat.Timeout,
		Stream:              stream,
	}
}

func (c *Session) Request(ctx context.Context) (*NonStreamingResponse, error) {
	if c.Stream {
		return nil, fmt.Errorf("stream request not supported")
	}
	body, err := jsoniter.Marshal(c)
	if err != nil {
		return nil, err
	}
	resp := new(NonStreamingResponse)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add(HeaderContentType, HeaderApplicationJson)
	req.Header.Add(HeaderAuthorization, fmt.Sprintf("Bearer %s", c.PersonalAccessToken))
	req.Header.Add(HeaderConnection, HeaderKeepAlive)
	req.Header.Add(HeaderAccept, HeaderAcceptAll)
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

func (c *Session) StreamRequest(ctx context.Context) (<-chan *StreamingResponse, <-chan error) {
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

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			errChan <- err
			return
		}
		req.Header.Add(HeaderContentType, HeaderApplicationJson)
		req.Header.Add(HeaderAuthorization, fmt.Sprintf("Bearer %s", c.PersonalAccessToken))
		req.Header.Add(HeaderConnection, HeaderKeepAlive)
		req.Header.Add(HeaderAccept, HeaderAcceptAll)

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
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "data:") {
				data := strings.TrimSpace(line[len("data:"):])
				if data == `{"event":"done"}` {
					break
				}
				var resp StreamingResponse
				if err = jsoniter.UnmarshalFromString(data, &resp); err != nil {
					errChan <- err
					return
				}
				respChan <- &resp
			} else if line != "" {
				var resp StreamingResponse
				if err = jsoniter.UnmarshalFromString(line, &resp); err != nil {
					errChan <- err
					return
				}
				respChan <- &resp
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

func (c *Session) WithConversationId(conversationId string) *Session {
	c.ConversationId = conversationId
	return c
}

func (c *Session) WithQuery(query string) *Session {
	c.Query = query
	return c
}

func (c *Session) WithChatHistory(chatHistories []Message) *Session {
	c.ChatHistory = chatHistories
	return c
}

func (c *Session) WithCustomVariables(customVariables map[string]any) *Session {
	c.CustomVariables = customVariables
	return c
}
