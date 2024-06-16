package coze

import "fmt"

type Message struct {
	// The role who returns the message.
	// 发送这条消息的实体。
	Role string `json:"role"`
	// The type of the message when the role is assistant.
	// 当 role = assistant 时，用于标识 Bot 的消息类型。
	Type string `json:"type"`
	// The returned content.
	// 消息内容。
	Content string `json:"content"`
	// The type of the return content.
	// 消息内容的类型。
	ContentType string `json:"content_type"`
}

type SessionResponse struct {
	// The ID of the conversation.
	// 会话 ID。
	ConversationId string `json:"conversation_id"`
	// The ID of the code.
	// 0 represents a successful call.
	// 状态码。
	// 0 代表调用成功。
	Code int `json:"code"`
	// The message of the request.
	// 状态信息。
	Msg string `json:"msg"`
}

type NonStreamingResponse struct {
	SessionResponse
	// The completed messages returned in JSON array.
	// 全部消息都处理完成后，以 JSON 数组形式返回。
	Messages []Message `json:"messages"`
}

type StreamingResponse struct {
	SessionResponse
	// The data set returned in the event.
	// 当前流式返回的数据包事件，不同 event 下，数据包会返回不同字段。
	Event string `json:"event"`
	// Whether the current message is completed.
	// 当前 message 是否结束。
	IsFinish bool `json:"is_finish"`
	// The identifier of the message. Each unique index corresponds to a single message.
	// 同一个 index 的增量返回属于同一条消息。
	Index int `json:"index"`
	// The incremental message that is returned at this time.
	// 增量返回的消息内容。
	Messages Message `json:"message"`
}

type HttpErrorResponse struct {
	Status     string // e.g. "200 OK"
	StatusCode int    `json:"status_code"` // http 状态码
	Body       []byte `json:"body"`        // http 响应体
}

func (h *HttpErrorResponse) Error() string {
	return fmt.Sprintf("response error: statusCode: %d, status: %s", h.StatusCode, h.Status)
}
