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

package conversation

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
	InternationalCreateUrl   = "https://api.coze.com/v1/conversation/create"
	InternationalRetrieveUrl = "https://api.coze.com/v1/conversation/retrieve"

	createUrl             = "https://api.coze.cn/v1/conversation/create"
	retrieveUrl           = "https://api.coze.cn/v1/conversation/retrieve"
	HeaderAuthorization   = "Authorization"
	HeaderContentType     = "Content-Type"
	HeaderApplicationJson = "application/json"
)

type Conversation struct {
	Authorization string                 `json:"-"`
	Timeout       time.Duration          `json:"-"`
	Messages      []request.EnterMessage `json:"messages,omitempty"`
	MetaData      map[string]any         `json:"meta_data,omitempty"`
}

func NewConversation(authorization string) *Conversation {
	return &Conversation{Authorization: authorization}
}

func (c *Conversation) WithMessages(messages ...request.EnterMessage) *Conversation {
	c.Messages = append(c.Messages, messages...)
	return c
}

func (c *Conversation) WithMetaData(metaData map[string]any) *Conversation {
	c.MetaData = metaData
	return c
}

func (c *Conversation) WithTimeout(timeout time.Duration) *Conversation {
	c.Timeout = timeout
	return c
}

func (c *Conversation) CreateRequest(ctx context.Context) (*response.DataResponse[response.Conversation], error) {
	body, err := jsoniter.Marshal(c)
	if err != nil {
		return nil, err
	}

	resp := new(response.DataResponse[response.Conversation])

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, createUrl, bytes.NewReader(body))
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

func (c *Conversation) RetrieveRequest(ctx context.Context, conversation string) (*response.DataResponse[response.Conversation], error) {
	resp := new(response.DataResponse[response.Conversation])

	// 构建查询参数
	params := url.Values{}
	params.Add("conversation_id", conversation)

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
