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
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/chenmingyong0423/go-coze/common/request"
	"github.com/stretchr/testify/require"

	"github.com/chenmingyong0423/go-coze/common/response"
)

func TestCreateRequest_Do(t *testing.T) {
	testCases := []struct {
		name           string
		ctx            context.Context
		authorization  string
		timeout        time.Duration
		conversationId string
		role           string
		content        string
		contentType    string
		objectString   request.ObjectString
		meta           map[string]any
		want           func(t *testing.T, resp *response.DataResponse[response.Message])
		wantErr        require.ErrorAssertionFunc
	}{
		{
			name:          "empty authorization",
			ctx:           context.Background(),
			authorization: "",
			wantErr: func(t require.TestingT, err error, i ...interface{}) {
				var errResp *response.HttpErrorResponse
				if errors.As(err, &errResp) {
					require.Equal(t, 401, errResp.StatusCode)
				} else {
					t.Errorf("unexpected error: %v", err)
				}
			},
		},
		{
			name:           "empty conversationId",
			ctx:            context.Background(),
			authorization:  os.Getenv("COZE_TOKEN"),
			conversationId: "",
			want: func(t *testing.T, resp *response.DataResponse[response.Message]) {
				require.Equal(t, 4002, resp.Code)
			},
		},
		{
			name:           "empty role",
			ctx:            context.Background(),
			authorization:  os.Getenv("COZE_TOKEN"),
			conversationId: "7414413032111063080",
			role:           "",
			content:        "你好",
			contentType:    "text",
			want: func(t *testing.T, resp *response.DataResponse[response.Message]) {
				require.Equal(t, 4000, resp.Code)
			},
		},
		{
			name:           "success",
			ctx:            context.Background(),
			authorization:  os.Getenv("COZE_TOKEN"),
			conversationId: "7414413032111063080",
			role:           "user",
			content:        "你好",
			contentType:    "text",
			want: func(t *testing.T, resp *response.DataResponse[response.Message]) {
				require.Equal(t, 0, resp.Code)
				require.NotNil(t, resp.Data)
				t.Log(resp.Data)
				require.NotZero(t, resp.Data.Id)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cr := NewMessage(tc.authorization, tc.conversationId).CreateRequest().WithRole(tc.role)
			if tc.contentType == "text" {
				cr.WithTextContent(tc.content)
			} else if tc.contentType == "object_string" {
				cr.WithObjectStringContent(tc.objectString)
			}

			got, err := cr.Do(tc.ctx)
			if err != nil {
				t.Log(err)
				require.NotNil(t, tc.wantErr)
				tc.wantErr(t, err)
			}
			if tc.want != nil {
				t.Log(got)
				tc.want(t, got)
			}
		})
	}
}

func TestListRequest_Do(t *testing.T) {
	testCases := []struct {
		name           string
		ctx            context.Context
		authorization  string
		timeout        time.Duration
		conversationId string
		order          string
		chatId         string
		beforeId       string
		afterId        string
		limit          int

		want    func(t *testing.T, resp *ListResponse[[]response.Message])
		wantErr require.ErrorAssertionFunc
	}{
		{
			name:           "empty authorization",
			ctx:            context.Background(),
			authorization:  "",
			conversationId: "7414413032111063080",
			wantErr: func(t require.TestingT, err error, i ...interface{}) {
				var errResp *response.HttpErrorResponse
				if errors.As(err, &errResp) {
					require.Equal(t, 401, errResp.StatusCode)
				} else {
					t.Errorf("unexpected error: %v", err)
				}
			},
		},
		{
			name:           "empty conversationId",
			ctx:            context.Background(),
			authorization:  os.Getenv("COZE_TOKEN"),
			conversationId: "",
			want: func(t *testing.T, resp *ListResponse[[]response.Message]) {
				require.Equal(t, 4002, resp.Code)
			},
		},
		{
			name:           "success",
			ctx:            context.Background(),
			authorization:  os.Getenv("COZE_TOKEN"),
			conversationId: "7414413032111063080",
			want: func(t *testing.T, resp *ListResponse[[]response.Message]) {
				require.Equal(t, 0, resp.Code)
				require.NotNil(t, resp.Data)
				t.Log(resp.Data)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lr := NewMessage(tc.authorization, tc.conversationId).ListRequest().
				WithOrder(tc.order).WithChatId(tc.chatId).WithBeforeId(tc.beforeId).WithAfterId(tc.afterId).WithLimit(tc.limit)

			got, err := lr.Do(tc.ctx)
			if err != nil {
				t.Log(err)
				require.NotNil(t, tc.wantErr)
				tc.wantErr(t, err)
			}
			if tc.want != nil {
				t.Log(got)
				tc.want(t, got)
			}
		})
	}
}

func TestRetrieveRequest_Do(t *testing.T) {
	resp, err := NewMessage(os.Getenv("COZE_TOKEN"), "7414413032111063080").ListRequest().Do(context.Background())
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotZero(t, len(resp.Data))
	messageId := resp.Data[0].Id

	testCases := []struct {
		name           string
		ctx            context.Context
		authorization  string
		timeout        time.Duration
		conversationId string
		messageId      string

		want    func(t *testing.T, resp *response.DataResponse[response.Message])
		wantErr require.ErrorAssertionFunc
	}{
		{
			name:           "empty authorization",
			ctx:            context.Background(),
			authorization:  "",
			conversationId: "7414413032111063080",
			messageId:      messageId,
			wantErr: func(t require.TestingT, err error, i ...interface{}) {
				var errResp *response.HttpErrorResponse
				if errors.As(err, &errResp) {
					require.Equal(t, 401, errResp.StatusCode)
				} else {
					t.Errorf("unexpected error: %v", err)
				}
			},
		},
		{
			name:           "empty conversationId",
			ctx:            context.Background(),
			authorization:  os.Getenv("COZE_TOKEN"),
			conversationId: "",
			messageId:      messageId,
			want: func(t *testing.T, resp *response.DataResponse[response.Message]) {
				require.Equal(t, 4002, resp.Code)
			},
		},
		{
			name:           "empty messageId",
			ctx:            context.Background(),
			authorization:  os.Getenv("COZE_TOKEN"),
			conversationId: "7414413032111063080",
			messageId:      "",
			want: func(t *testing.T, resp *response.DataResponse[response.Message]) {
				require.Equal(t, 4005, resp.Code)
			},
		},
		{
			name:           "success",
			ctx:            context.Background(),
			authorization:  os.Getenv("COZE_TOKEN"),
			conversationId: "7414413032111063080",
			messageId:      messageId,
			want: func(t *testing.T, resp *response.DataResponse[response.Message]) {
				require.Equal(t, 0, resp.Code)
				require.NotNil(t, resp.Data)
				t.Log(resp.Data)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := NewMessage(tc.authorization, tc.conversationId).RetrieveRequest(tc.messageId)

			got, err := rr.Do(tc.ctx)
			if err != nil {
				t.Log(err)
				require.NotNil(t, tc.wantErr)
				tc.wantErr(t, err)
			}
			if tc.want != nil {
				t.Log(got)
				tc.want(t, got)
			}
		})
	}
}

func TestModifyRequest_Do(t *testing.T) {
	resp, err := NewMessage(os.Getenv("COZE_TOKEN"), "7414413032111063080").ListRequest().Do(context.Background())
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotZero(t, len(resp.Data))
	messageId := resp.Data[0].Id

	testCases := []struct {
		name           string
		ctx            context.Context
		authorization  string
		timeout        time.Duration
		conversationId string
		messageId      string
		content        string
		contentType    string
		objectString   request.ObjectString
		meta           map[string]any

		want    func(t *testing.T, resp *ModifyResponse[response.Message])
		wantErr require.ErrorAssertionFunc
	}{
		{
			name:           "empty authorization",
			ctx:            context.Background(),
			authorization:  "",
			conversationId: "7414413032111063080",
			messageId:      messageId,
			wantErr: func(t require.TestingT, err error, i ...interface{}) {
				var errResp *response.HttpErrorResponse
				if errors.As(err, &errResp) {
					require.Equal(t, 401, errResp.StatusCode)
				} else {
					t.Errorf("unexpected error: %v", err)
				}
			},
		},
		{
			name:           "empty conversationId",
			ctx:            context.Background(),
			authorization:  os.Getenv("COZE_TOKEN"),
			conversationId: "",
			messageId:      messageId,
			want: func(t *testing.T, resp *ModifyResponse[response.Message]) {
				require.Equal(t, 4002, resp.Code)
			},
		},
		{
			name:           "empty messageId",
			ctx:            context.Background(),
			authorization:  os.Getenv("COZE_TOKEN"),
			conversationId: "7414413032111063080",
			messageId:      "",
			want: func(t *testing.T, resp *ModifyResponse[response.Message]) {
				require.Equal(t, 4005, resp.Code)
			},
		},
		{
			name:           "success",
			ctx:            context.Background(),
			authorization:  os.Getenv("COZE_TOKEN"),
			conversationId: "7414413032111063080",
			messageId:      messageId,
			content:        "test",
			contentType:    "text",
			want: func(t *testing.T, resp *ModifyResponse[response.Message]) {
				require.Equal(t, 0, resp.Code)
				require.NotNil(t, resp.Message)
				t.Log(resp.Message)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := NewMessage(tc.authorization, tc.conversationId).ModifyRequest(tc.messageId)

			if tc.contentType == "text" {
				rr.WithTextContent(tc.content)
			} else if tc.contentType == "object_string" {
				rr.WithObjectStringContent(tc.objectString)
			}

			got, err := rr.Do(tc.ctx)
			if err != nil {
				t.Log(err)
				require.NotNil(t, tc.wantErr)
				tc.wantErr(t, err)
			}
			if tc.want != nil {
				t.Log(got)
				tc.want(t, got)
			}
		})
	}
}

func TestDeleteRequest_Do(t *testing.T) {
	resp, err := NewMessage(os.Getenv("COZE_TOKEN"), "7414413032111063080").ListRequest().Do(context.Background())
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotZero(t, len(resp.Data))
	messageId := resp.Data[0].Id

	testCases := []struct {
		name           string
		ctx            context.Context
		authorization  string
		timeout        time.Duration
		conversationId string
		messageId      string

		want    func(t *testing.T, resp *response.DataResponse[response.Message])
		wantErr require.ErrorAssertionFunc
	}{
		{
			name:           "empty authorization",
			ctx:            context.Background(),
			authorization:  "",
			conversationId: "7414413032111063080",
			messageId:      messageId,
			wantErr: func(t require.TestingT, err error, i ...interface{}) {
				var errResp *response.HttpErrorResponse
				if errors.As(err, &errResp) {
					require.Equal(t, 401, errResp.StatusCode)
				} else {
					t.Errorf("unexpected error: %v", err)
				}
			},
		},
		{
			name:           "empty conversationId",
			ctx:            context.Background(),
			authorization:  os.Getenv("COZE_TOKEN"),
			conversationId: "",
			messageId:      messageId,
			want: func(t *testing.T, resp *response.DataResponse[response.Message]) {
				require.Equal(t, 4002, resp.Code)
			},
		},
		{
			name:           "empty messageId",
			ctx:            context.Background(),
			authorization:  os.Getenv("COZE_TOKEN"),
			conversationId: "7414413032111063080",
			messageId:      "",
			want: func(t *testing.T, resp *response.DataResponse[response.Message]) {
				require.Equal(t, 4005, resp.Code)
			},
		},
		{
			name:           "success",
			ctx:            context.Background(),
			authorization:  os.Getenv("COZE_TOKEN"),
			conversationId: "7414413032111063080",
			messageId:      messageId,
			want: func(t *testing.T, resp *response.DataResponse[response.Message]) {
				require.Equal(t, 0, resp.Code)
				require.NotNil(t, resp.Data)
				t.Log(resp.Data)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dr := NewMessage(tc.authorization, tc.conversationId).DeleteRequest(tc.messageId)

			got, err := dr.Do(tc.ctx)
			if err != nil {
				t.Log(err)
				require.NotNil(t, tc.wantErr)
				tc.wantErr(t, err)
			}
			if tc.want != nil {
				t.Log(got)
				tc.want(t, got)
			}
		})
	}
}
