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
	HeaderAuthorization   = "authorization"
	HeaderContentType     = "Content-Type"
	HeaderApplicationJson = "application/json"
)

type Conversation struct {
	authorization string
}

func NewConversation(authorization string) *Conversation {
	return &Conversation{authorization: authorization}
}

func (c *Conversation) CreateRequest() *CreateRequest {
	return &CreateRequest{
		conversation: c,
	}
}

func (c *Conversation) RetrieveRequest() *RetrieveRequest {
	return &RetrieveRequest{
		conversation: c,
	}
}

type CreateRequest struct {
	conversation *Conversation

	timeout  time.Duration
	Messages []request.EnterMessage `json:"messages,omitempty"`
	MetaData map[string]any         `json:"meta_data,omitempty"`
}

func (r *CreateRequest) WithMessages(messages ...request.EnterMessage) *CreateRequest {
	r.Messages = append(r.Messages, messages...)
	return r
}

func (r *CreateRequest) WithMetaData(metaData map[string]any) *CreateRequest {
	r.MetaData = metaData
	return r
}

func (r *CreateRequest) WithTimeout(timeout time.Duration) *CreateRequest {
	r.timeout = timeout
	return r
}

func (r *CreateRequest) Do(ctx context.Context) (*response.DataResponse[response.Conversation], error) {
	body, err := jsoniter.Marshal(r)
	if err != nil {
		return nil, err
	}

	resp := new(response.DataResponse[response.Conversation])

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, createUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add(HeaderContentType, HeaderApplicationJson)
	req.Header.Add(HeaderAuthorization, fmt.Sprintf("Bearer %s", r.conversation.authorization))

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

type RetrieveRequest struct {
	conversation *Conversation

	timeout time.Duration
}

func (r *RetrieveRequest) WithTimeout(timeout time.Duration) *RetrieveRequest {
	r.timeout = timeout
	return r
}

func (r *RetrieveRequest) Do(ctx context.Context, conversation string) (*response.DataResponse[response.Conversation], error) {
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
	req.Header.Add(HeaderAuthorization, fmt.Sprintf("Bearer %s", r.conversation.authorization))

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
