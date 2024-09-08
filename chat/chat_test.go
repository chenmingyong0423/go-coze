package chat

import (
	"context"
	"os"
	"testing"

	"github.com/chenmingyong0423/go-coze/common/request"

	"github.com/stretchr/testify/require"
)

func TestSession_Request(t *testing.T) {

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

		wantError require.ErrorAssertionFunc
		wantCode  int
	}{
		{
			name:               "invalid request params for botID",
			ctx:                context.Background(),
			botID:              "",
			userId:             os.Getenv("COZE_USER_ID"),
			authorization:      os.Getenv("COZE_TOKEN"),
			stream:             false,
			conversationId:     "",
			additionalMessages: nil,
			customVariables:    nil,
			autoSaveHistory:    false,
			metaData:           nil,
			extraParams:        nil,
			wantError:          require.NoError,
			wantCode:           4006,
		},
		{
			name:           "invalid request params for userId",
			ctx:            context.Background(),
			botID:          os.Getenv("COZE_BOT_ID"),
			userId:         "",
			authorization:  os.Getenv("COZE_TOKEN"),
			stream:         false,
			conversationId: "",
			additionalMessages: []request.EnterMessage{
				request.NewEnterMessageBuilder().Role("user").Content("你好").ContentType("text").Build(),
			},
			customVariables: nil,
			autoSaveHistory: false,
			metaData:        nil,
			extraParams:     nil,
			wantError:       require.NoError,
			wantCode:        4000,
		},
		{
			name:           "invalid request params for token",
			ctx:            context.Background(),
			botID:          os.Getenv("COZE_BOT_ID"),
			userId:         os.Getenv("COZE_USER_ID"),
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
			wantError:       require.NoError,
			wantCode:        4100,
		},
		{
			name:               "empty content",
			ctx:                context.Background(),
			botID:              os.Getenv("COZE_BOT_ID"),
			userId:             os.Getenv("COZE_USER_ID"),
			authorization:      os.Getenv("COZE_TOKEN"),
			stream:             false,
			conversationId:     "",
			additionalMessages: nil,
			customVariables:    nil,
			autoSaveHistory:    false,
			metaData:           nil,
			extraParams:        nil,
			wantError:          require.NoError,
			wantCode:           0,
		},
		{
			name:           "empty content",
			ctx:            context.Background(),
			botID:          os.Getenv("COZE_BOT_ID"),
			userId:         os.Getenv("COZE_USER_ID"),
			authorization:  os.Getenv("COZE_TOKEN"),
			stream:         false,
			conversationId: "",
			additionalMessages: []request.EnterMessage{
				request.NewEnterMessageBuilder().Role("user").Content("你好").ContentType("text").Build(),
			},
			customVariables: nil,
			autoSaveHistory: false,
			metaData:        nil,
			extraParams:     nil,
			wantError:       require.NoError,
			wantCode:        0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建一个聊天对象
			chat := NewChat(tc.authorization, tc.userId, tc.botID)

			// 添加请求参数并发送以及处理错误
			resp, err := chat.WithConversationId(tc.conversationId).
				AddMessages(tc.additionalMessages...).
				WithCustomVariables(tc.customVariables).
				WithAutoSaveHistory(tc.autoSaveHistory).
				WithMetaData(tc.metaData).
				WithExtraParams(tc.extraParams).
				Chat(tc.ctx)
			tc.wantError(t, err)
			if err == nil {
				t.Log(resp)
				if resp.Data != nil {
					t.Log(resp.Data)
				}
				require.Equal(t, tc.wantCode, resp.Code)
			}
		})
	}
}

func TestSession_StreamRequest(t *testing.T) {
	{
		// 创建一个聊天对象
		chat := NewChat("", "", "")

		respChan, errChan := chat.WithStream(true).WithConversationId("").
			AddMessages(request.NewEnterMessageBuilder().Role("user").Content("你好").ContentType("text").Build()).
			WithCustomVariables(nil).
			WithAutoSaveHistory(false).
			WithMetaData(nil).
			WithExtraParams(nil).
			StreamChat(context.Background())
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
		// 正常请求
		botId := os.Getenv("COZE_BOT_ID")
		userId := os.Getenv("COZE_USER_ID")
		token := os.Getenv("COZE_TOKEN")

		chat := NewChat(token, userId, botId)
		respChan, errChan := chat.WithStream(true).WithConversationId("").
			AddMessages(request.NewEnterMessageBuilder().Role("user").Content("你好").ContentType("text").Build()).
			WithCustomVariables(nil).
			WithAutoSaveHistory(false).
			WithMetaData(nil).
			WithExtraParams(nil).
			StreamChat(context.Background())
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
