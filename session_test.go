package coze

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSession_Request(t *testing.T) {

	testCases := []struct {
		name                string
		ctx                 context.Context
		botID               string
		user                string
		personalAccessToken string

		stream          bool
		conversationId  string
		query           string
		chatHistory     []Message
		customVariables map[string]any

		wantError require.ErrorAssertionFunc
		wantCode  int
	}{
		{
			name:                "invalid request params for botID",
			ctx:                 context.Background(),
			botID:               "",
			user:                os.Getenv("COZE_USER"),
			personalAccessToken: os.Getenv("COZE_TOKEN"),
			stream:              false,
			query:               "你好",
			chatHistory:         nil,
			customVariables:     nil,
			wantCode:            702242001,
			wantError:           require.NoError,
		},
		{
			name:                "invalid request params for user",
			ctx:                 context.Background(),
			botID:               os.Getenv("COZE_BOT_ID"),
			user:                "",
			personalAccessToken: os.Getenv("COZE_TOKEN"),
			stream:              false,
			query:               "你好",
			chatHistory:         nil,
			customVariables:     nil,
			wantCode:            702242001,
			wantError:           require.NoError,
		},
		{
			name:                "invalid request params for token",
			ctx:                 context.Background(),
			botID:               os.Getenv("COZE_BOT_ID"),
			user:                os.Getenv("COZE_USER"),
			personalAccessToken: "",
			stream:              false,
			query:               "你好",
			chatHistory:         nil,
			customVariables:     nil,
			wantCode:            702242001,
			wantError:           require.NoError,
		},

		{
			name:                "invalid request params for query",
			ctx:                 context.Background(),
			botID:               os.Getenv("COZE_BOT_ID"),
			user:                os.Getenv("COZE_USER"),
			personalAccessToken: os.Getenv("COZE_TOKEN"),
			stream:              false,
			query:               "",
			chatHistory:         nil,
			customVariables:     nil,
			wantCode:            702242001,
			wantError:           require.NoError,
		},
		{
			name:                "success",
			ctx:                 context.Background(),
			botID:               os.Getenv("COZE_BOT_ID"),
			user:                os.Getenv("COZE_USER"),
			personalAccessToken: os.Getenv("COZE_TOKEN"),
			stream:              false,
			query:               "你好",
			chatHistory:         nil,
			customVariables:     nil,
			wantError:           require.NoError,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建一个聊天对象
			chat := NewChat(tc.botID, tc.user, tc.personalAccessToken)

			// 创建新的会话对象并设置会话流和类型
			session := chat.Chat(tc.stream)

			// 添加请求参数并发送以及处理错误
			resp, err := session.WithQuery(tc.query).
				WithConversationId(tc.conversationId).
				WithChatHistory(tc.chatHistory).
				WithCustomVariables(tc.customVariables).
				Request(tc.ctx)
			tc.wantError(t, err)
			if err == nil {
				t.Log(resp)
				require.Equal(t, tc.wantCode, resp.Code)
			}
		})
	}
}

func TestSession_StreamRequest(t *testing.T) {
	{
		// 创建一个聊天对象
		chat := NewChat("", "", "")
		session := chat.Chat(true)
		respChan, errChan := session.WithQuery("你好").
			WithConversationId("").
			WithChatHistory(nil).
			WithCustomVariables(nil).
			StreamRequest(context.Background())
		for {
			select {
			case resp, ok := <-respChan:
				require.True(t, ok)
				require.NotNil(t, resp)
				t.Log(resp)
				require.Equal(t, resp.Code, 702242001)
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

		user := os.Getenv("COZE_USER")

		token := os.Getenv("COZE_TOKEN")

		chat := NewChat(botId, user, token)
		session := chat.Chat(true)
		respChan, errChan := session.WithQuery("你好").
			WithConversationId("").
			WithChatHistory(nil).
			WithCustomVariables(nil).
			StreamRequest(context.Background())
		for {
			select {
			case resp, ok := <-respChan:
				if !ok {
					respChan = nil
				} else {
					require.NotNil(t, resp)
					require.Equal(t, resp.Code, 0)
					t.Log(resp)
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
