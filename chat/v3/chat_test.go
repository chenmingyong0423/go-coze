package chat

import (
	"context"
	"testing"
	"time"

	"github.com/chenmingyong0423/go-coze/common/response"

	"github.com/chenmingyong0423/go-coze/common/request"

	"github.com/stretchr/testify/require"
)

func TestCreateRequest_Do(t *testing.T) {
	testCases := []struct {
		name          string
		ctx           context.Context
		botID         string
		userId        string
		authorization string

		stream             bool
		conversationId     string
		additionalMessages []request.EnterMessage
		customVariables    map[string]any
		autoSaveHistory    bool
		metaData           map[string]any
		extraParams        []string

		wantErrorFunc require.ErrorAssertionFunc
		wantFunc      func(t *testing.T, resp *response.DataResponse[*response.Chat])
	}{
		{
			name:               "invalid request params for botID",
			ctx:                context.Background(),
			botID:              "",
			userId:             "7330571112050343973",
			authorization:      "pat_MPmFilLOZIU1VsXa6bC4mrUCyFIULmaxpYIaRRWe1I77n96dLIVfwW5ucGKt5kqP",
			stream:             false,
			conversationId:     "",
			additionalMessages: nil,
			customVariables:    nil,
			autoSaveHistory:    false,
			metaData:           nil,
			extraParams:        nil,
			wantFunc: func(t *testing.T, resp *response.DataResponse[*response.Chat]) {
				require.NotNil(t, resp)
				require.Equal(t, 4006, resp.Code)
			},
		},
		{
			name:           "invalid request params for userId",
			ctx:            context.Background(),
			botID:          "7378912442585874447",
			userId:         "",
			authorization:  "pat_MPmFilLOZIU1VsXa6bC4mrUCyFIULmaxpYIaRRWe1I77n96dLIVfwW5ucGKt5kqP",
			stream:         false,
			conversationId: "",
			additionalMessages: []request.EnterMessage{
				request.NewEnterMessageBuilder().Role("user").Content("你好").ContentType("text").Build(),
			},
			customVariables: nil,
			autoSaveHistory: false,
			metaData:        nil,
			extraParams:     nil,
			wantFunc: func(t *testing.T, resp *response.DataResponse[*response.Chat]) {
				require.NotNil(t, resp)
				require.Equal(t, 4000, resp.Code)
			},
		},
		{
			name:           "invalid request params for token",
			ctx:            context.Background(),
			botID:          "7378912442585874447",
			userId:         "7330571112050343973",
			authorization:  "",
			stream:         false,
			conversationId: "",
			additionalMessages: []request.EnterMessage{
				request.NewEnterMessageBuilder().Role("user").Content("你好").ContentType("text").Build(),
			},
			customVariables: nil,
			autoSaveHistory: false,
			metaData:        nil,
			extraParams:     nil,
			wantFunc: func(t *testing.T, resp *response.DataResponse[*response.Chat]) {
				require.NotNil(t, resp)
				require.Equal(t, 4100, resp.Code)
			},
		},
		{
			name:               "empty content",
			ctx:                context.Background(),
			botID:              "7378912442585874447",
			userId:             "7330571112050343973",
			authorization:      "pat_MPmFilLOZIU1VsXa6bC4mrUCyFIULmaxpYIaRRWe1I77n96dLIVfwW5ucGKt5kqP",
			stream:             false,
			conversationId:     "",
			additionalMessages: nil,
			customVariables:    nil,
			autoSaveHistory:    false,
			metaData:           nil,
			extraParams:        nil,
			wantErrorFunc:      require.NoError,
			wantFunc: func(t *testing.T, resp *response.DataResponse[*response.Chat]) {
				require.NotNil(t, resp)
				require.Equal(t, 0, resp.Code)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			chat := NewChat(tc.authorization, tc.userId, tc.botID)

			got, err := chat.ChatRequest().WithConversationId(tc.conversationId).
				AddMessages(tc.additionalMessages...).
				WithCustomVariables(tc.customVariables).
				WithAutoSaveHistory(tc.autoSaveHistory).
				WithMetaData(tc.metaData).
				WithExtraParams(tc.extraParams).
				Do(tc.ctx)
			if err != nil {
				t.Log(err)
				require.NotNil(t, tc.wantErrorFunc)
				tc.wantErrorFunc(t, err)
			}
			if tc.wantFunc != nil {
				t.Log(got)
				tc.wantFunc(t, got)
			}
		})
	}
}

func TestCreateRequest_DoStream(t *testing.T) {
	{
		chat := NewChat("", "", "")

		respChan, errChan := chat.ChatRequest().WithConversationId("").
			AddMessages(request.NewEnterMessageBuilder().Role("user").Content("你好").ContentType("text").Build()).
			WithCustomVariables(nil).
			WithAutoSaveHistory(false).
			WithMetaData(nil).
			WithExtraParams(nil).
			DoStream(context.Background())
		for {
			select {
			case resp, ok := <-respChan:
				require.True(t, ok)
				require.NotNil(t, resp)
				t.Log(resp)
				require.Equal(t, 4100, resp.Code)
				respChan = nil
			case err, ok := <-errChan:
				require.False(t, ok)
				require.Nil(t, err)
				errChan = nil
			}
			if respChan == nil && errChan == nil {
				break
			}
		}
	}
	{
		botId := "7378912442585874447"
		userId := "7330571112050343973"
		token := "pat_MPmFilLOZIU1VsXa6bC4mrUCyFIULmaxpYIaRRWe1I77n96dLIVfwW5ucGKt5kqP"

		chat := NewChat(token, userId, botId)
		respChan, errChan := chat.ChatRequest().WithConversationId("").
			AddMessages(request.NewEnterMessageBuilder().Role("user").Content("你好").ContentType("text").Build()).
			WithCustomVariables(nil).
			WithAutoSaveHistory(false).
			WithMetaData(nil).
			WithExtraParams(nil).
			DoStream(context.Background())
		for {
			select {
			case resp, ok := <-respChan:
				if !ok {
					respChan = nil
				} else {
					require.True(t, ok)
					t.Log(resp)
					require.Zero(t, resp.Code)
					if resp.Chat != nil {
						t.Log(resp.Chat)
					}
					if resp.Message != nil {
						t.Log(resp.Message)
					}
				}
			case err, ok := <-errChan:
				require.False(t, ok)
				require.Nil(t, err)
				errChan = nil
			}
			if respChan == nil && errChan == nil {
				break
			}
		}
	}
}

func TestRetrieveRequest_Do(t *testing.T) {
	chat := NewChat("pat_MPmFilLOZIU1VsXa6bC4mrUCyFIULmaxpYIaRRWe1I77n96dLIVfwW5ucGKt5kqP", "7330571112050343973", "7378912442585874447")
	resp, err := chat.ChatRequest().WithAutoSaveHistory(true).AddMessages(request.NewEnterMessageBuilder().Role("user").Content("你好").ContentType("text").Build()).
		Do(context.Background())
	require.NoError(t, err)
	t.Log(resp)
	require.Equal(t, 0, resp.Code)
	resp2, err := chat.RetrieveRequest(resp.Data.ConversationId).Do(context.Background(), resp.Data.Id)
	require.NoError(t, err)
	require.Equal(t, 0, resp.Code)
	require.NotNil(t, resp2)
	t.Log(resp2.Data)
}

func TestMessageListRequest_Do(t *testing.T) {
	chat := NewChat("pat_MPmFilLOZIU1VsXa6bC4mrUCyFIULmaxpYIaRRWe1I77n96dLIVfwW5ucGKt5kqP", "7330571112050343973", "7378912442585874447")
	resp, err := chat.ChatRequest().WithAutoSaveHistory(true).AddMessages(request.NewEnterMessageBuilder().Role("user").Content("你好").ContentType("text").Build()).
		Do(context.Background())
	require.NoError(t, err)
	t.Log(resp)
	require.Equal(t, 0, resp.Code)
	// 等待对话完成
	time.Sleep(2 * time.Second)

	resp2, err := chat.MessageListRequest(resp.Data.ConversationId).Do(context.Background(), resp.Data.Id)
	require.NoError(t, err)
	t.Log(resp2)
	require.Equal(t, 0, resp2.Code)
	for _, message := range resp2.Data {
		t.Log(message)
	}
}

func TestCancelRequest_Do(t *testing.T) {
	chat := NewChat("pat_MPmFilLOZIU1VsXa6bC4mrUCyFIULmaxpYIaRRWe1I77n96dLIVfwW5ucGKt5kqP", "7330571112050343973", "7378912442585874447")
	resp, err := chat.ChatRequest().WithAutoSaveHistory(true).AddMessages(request.NewEnterMessageBuilder().Role("user").Content("你好").ContentType("text").Build()).
		Do(context.Background())
	require.NoError(t, err)
	t.Log(resp)
	require.Equal(t, 0, resp.Code)

	resp2, err := chat.CancelRequest(resp.Data.ConversationId).Do(context.Background(), resp.Data.Id)
	require.NoError(t, err)
	require.Equal(t, 0, resp.Code)
	require.NotNil(t, resp2)
	t.Log(resp2.Data)
}
