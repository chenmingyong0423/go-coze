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
	"context"
	"errors"
	"github.com/chenmingyong0423/go-coze/common/request"
	"github.com/chenmingyong0423/go-coze/common/response"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestConversation_CreateRequest(t *testing.T) {

	tests := []struct {
		name          string
		ctx           context.Context
		authorization string
		timeout       time.Duration
		messages      []request.EnterMessage
		metaData      map[string]any
		want          func(t *testing.T, resp *response.DataResponse[response.Conversation])
		wantErr       require.ErrorAssertionFunc
	}{
		{
			name:          "empty authorization",
			ctx:           context.Background(),
			authorization: "",
			metaData:      map[string]any{},
			wantErr: func(t require.TestingT, err error, i ...interface{}) {
				var errResp *response.HttpErrorResponse
				if errors.As(err, &errResp) {
					require.Equal(t, 401, errResp.StatusCode)
				} else {
					t.FailNow()
				}
			},
		},
		{
			name:          "success",
			ctx:           context.Background(),
			authorization: os.Getenv("COZE_TOKEN"),
			want: func(t *testing.T, resp *response.DataResponse[response.Conversation]) {
				require.Equal(t, 0, resp.Code)
				require.NotNil(t, resp.Data)
				t.Log(resp.Data)
				require.NotZero(t, resp.Data.Id)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConversation(tt.authorization).WithMetaData(tt.metaData).WithTimeout(tt.timeout).WithMessages(tt.messages...)
			got, err := c.CreateRequest(tt.ctx)
			if err != nil {
				t.Log(err)
				require.NotNil(t, tt.wantErr)
				tt.wantErr(t, err)
			}
			if tt.want != nil {
				tt.want(t, got)
			}
		})
	}
}

func TestConversation_RetrieveRequest(t *testing.T) {
	conversation := NewConversation(os.Getenv("COZE_TOKEN"))
	createResp, err := conversation.CreateRequest(context.Background())
	require.NoError(t, err)
	require.NotNil(t, createResp.Data)
	retrieveResp, err := conversation.RetrieveRequest(context.Background(), createResp.Data.Id)
	require.NoError(t, err)
	require.Equal(t, retrieveResp.Data, createResp.Data)
}
