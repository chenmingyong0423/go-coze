// Copyright 2023 chenmingyong0423

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package message

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/chenmingyong0423/go-coze/common/request"

	"github.com/chenmingyong0423/go-coze/common/response"
	jsoniter "github.com/json-iterator/go"
)

const (
	InternationalCreateUrl   = "https://api.coze.com/v1/conversation/message/create"
	InternationalListUrl     = "https://api.coze.cn/v1/conversation/message/list"
	InternationalRetrieveUrl = "https://api.coze.com/v1/conversation/message/retrieve"
	InternationalModifyUrl   = "https://api.coze.cn/v1/conversation/message/modify"
	InternationalDeleteUrl   = "https://api.coze.cn/v1/conversation/message/delete"

	createUrl             = "https://api.coze.cn/v1/conversation/message/create"
	listUrl               = "https://api.coze.cn/v1/conversation/message/list"
	retrieveUrl           = "https://api.coze.cn/v1/conversation/message/retrieve"
	modifyUrl             = "https://api.coze.cn/v1/conversation/message/modify"
	deleteUrl             = "https://api.coze.cn/v1/conversation/message/delete"
	HeaderAuthorization   = "Authorization"
	HeaderContentType     = "Content-Type"
	HeaderApplicationJson = "application/json"
)

type Message struct {
	authorization  string
	conversationId string
}

func NewMessage(authorization string, conversationId string) *Message {
	return &Message{authorization: authorization, conversationId: conversationId}
}

func (m *Message) CreateRequest() *CreateRequest {
	return &CreateRequest{message: m}
}

func (m *Message) ListRequest() *ListRequest {
	return &ListRequest{message: m}
}

func (m *Message) RetrieveRequest(messageId string) *RetrieveRequest {
	return &RetrieveRequest{message: m, messageId: messageId}
}

func (m *Message) ModifyRequest(messageId string) *ModifyRequest {
	return &ModifyRequest{message: m, messageId: messageId}
}

func (m *Message) DeleteRequest(messageId string) *DeleteRequest {
	return &DeleteRequest{message: m, messageId: messageId}
}

type CreateRequest struct {
	timeout time.Duration

	message     *Message
	Role        string         `json:"role"`
	Content     string         `json:"content"`
	ContentType string         `json:"content_type"`
	Meta        map[string]any `json:"meta_data,omitempty"`
}

func (c *CreateRequest) WithTimeout(timeout time.Duration) *CreateRequest {
	c.timeout = timeout
	return c
}

func (c *CreateRequest) WithRole(role string) *CreateRequest {
	c.Role = role
	return c
}

func (c *CreateRequest) WithTextContent(content string) *CreateRequest {
	c.Content = content
	c.ContentType = "text"
	return c
}

func (c *CreateRequest) WithObjectStringContent(objectString request.ObjectString) *CreateRequest {
	marshal, _ := jsoniter.Marshal(objectString)
	c.Content = string(marshal)
	c.ContentType = "object_string"
	return c
}

func (c *CreateRequest) Do(ctx context.Context) (*response.DataResponse[response.Message], error) {
	params := url.Values{}
	params.Add("conversation_id", c.message.conversationId)

	u, err := url.Parse(createUrl)
	if err != nil {
		return nil, err
	}
	u.RawQuery = params.Encode()

	body, err := jsoniter.Marshal(c)
	if err != nil {
		return nil, err
	}

	resp := new(response.DataResponse[response.Message])

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add(HeaderContentType, HeaderApplicationJson)
	req.Header.Add(HeaderAuthorization, fmt.Sprintf("Bearer %s", c.message.authorization))

	client := http.DefaultClient
	if c.timeout != 0 {
		client.Timeout = c.timeout
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

type ListRequest struct {
	timeout time.Duration
	message *Message

	Order    string `json:"order,omitempty"`
	ChatId   string `json:"chat_id,omitempty"`
	BeforeId string `json:"before_id,omitempty"`
	AfterId  string `json:"after_id,omitempty"`
	Limit    int    `json:"limit,omitempty"`
}

func (c *ListRequest) WithTimeout(timeout time.Duration) *ListRequest {
	c.timeout = timeout
	return c
}

func (c *ListRequest) WithOrder(order string) *ListRequest {
	c.Order = order
	return c
}

func (c *ListRequest) WithChatId(chatId string) *ListRequest {
	c.ChatId = chatId
	return c
}

func (c *ListRequest) WithBeforeId(beforeId string) *ListRequest {
	c.BeforeId = beforeId
	return c
}

func (c *ListRequest) WithAfterId(afterId string) *ListRequest {
	c.AfterId = afterId
	return c
}

func (c *ListRequest) WithLimit(limit int) *ListRequest {
	c.Limit = limit
	return c
}

func (c *ListRequest) Do(ctx context.Context) (*ListResponse[[]response.Message], error) {
	params := url.Values{}
	params.Add("conversation_id", c.message.conversationId)

	u, err := url.Parse(listUrl)
	if err != nil {
		return nil, err
	}
	u.RawQuery = params.Encode()

	body, err := jsoniter.Marshal(c)
	if err != nil {
		return nil, err
	}

	resp := new(ListResponse[[]response.Message])
	resp.Data = make([]response.Message, 0)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add(HeaderContentType, HeaderApplicationJson)
	req.Header.Add(HeaderAuthorization, fmt.Sprintf("Bearer %s", c.message.authorization))

	client := http.DefaultClient
	if c.timeout != 0 {
		client.Timeout = c.timeout
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

type RetrieveRequest struct {
	timeout time.Duration
	message *Message

	messageId string
}

func (c *RetrieveRequest) WithTimeout(timeout time.Duration) *RetrieveRequest {
	c.timeout = timeout
	return c
}

func (c *RetrieveRequest) Do(ctx context.Context) (*response.DataResponse[response.Message], error) {
	params := url.Values{}
	params.Add("conversation_id", c.message.conversationId)
	params.Add("message_id", c.messageId)

	u, err := url.Parse(retrieveUrl)
	if err != nil {
		return nil, err
	}
	u.RawQuery = params.Encode()

	resp := new(response.DataResponse[response.Message])

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add(HeaderContentType, HeaderApplicationJson)
	req.Header.Add(HeaderAuthorization, fmt.Sprintf("Bearer %s", c.message.authorization))

	client := http.DefaultClient
	if c.timeout != 0 {
		client.Timeout = c.timeout
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

type ModifyRequest struct {
	timeout   time.Duration
	message   *Message
	messageId string

	Content     string         `json:"content,omitempty"`
	ContentType string         `json:"content_type,omitempty"`
	Meta        map[string]any `json:"meta_data,omitempty"`
}

func (c *ModifyRequest) WithTimeout(timeout time.Duration) *ModifyRequest {
	c.timeout = timeout
	return c
}

func (c *ModifyRequest) WithTextContent(content string) *ModifyRequest {
	c.Content = content
	c.ContentType = "text"
	return c
}

func (c *ModifyRequest) WithObjectStringContent(objectString request.ObjectString) *ModifyRequest {
	marshal, _ := jsoniter.Marshal(objectString)
	c.Content = string(marshal)
	c.ContentType = "object_string"
	return c
}

func (c *ModifyRequest) Do(ctx context.Context) (*ModifyResponse[response.Message], error) {
	params := url.Values{}
	params.Add("conversation_id", c.message.conversationId)
	params.Add("message_id", c.messageId)

	u, err := url.Parse(modifyUrl)
	if err != nil {
		return nil, err
	}
	u.RawQuery = params.Encode()

	body, err := jsoniter.Marshal(c)
	if err != nil {
		return nil, err
	}

	resp := new(ModifyResponse[response.Message])

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add(HeaderContentType, HeaderApplicationJson)
	req.Header.Add(HeaderAuthorization, fmt.Sprintf("Bearer %s", c.message.authorization))

	client := http.DefaultClient
	if c.timeout != 0 {
		client.Timeout = c.timeout
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

type DeleteRequest struct {
	timeout   time.Duration
	message   *Message
	messageId string
}

func (c *DeleteRequest) WithTimeout(timeout time.Duration) *DeleteRequest {
	c.timeout = timeout
	return c
}

func (c *DeleteRequest) Do(ctx context.Context) (*response.DataResponse[response.Message], error) {
	params := url.Values{}
	params.Add("conversation_id", c.message.conversationId)
	params.Add("message_id", c.messageId)

	u, err := url.Parse(deleteUrl)
	if err != nil {
		return nil, err
	}
	u.RawQuery = params.Encode()

	resp := new(response.DataResponse[response.Message])

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add(HeaderContentType, HeaderApplicationJson)
	req.Header.Add(HeaderAuthorization, fmt.Sprintf("Bearer %s", c.message.authorization))

	client := http.DefaultClient
	if c.timeout != 0 {
		client.Timeout = c.timeout
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
