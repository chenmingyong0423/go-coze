package chat

import (
	"github.com/chenmingyong0423/go-coze/common/response"
)

type StreamingResponse struct {
	response.BaseResponse
	Event   string
	Chat    *response.Chat
	Message *response.Message
}
