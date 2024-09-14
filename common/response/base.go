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

package response

import "fmt"

type BaseResponse struct {
	// The ID of the code.
	// 0 represents a successful call.
	// 状态码。
	// 0 代表调用成功。
	Code int `json:"code"`
	// 状态信息。API 调用失败时可通过此字段查看详细错误信息。
	Msg string `json:"msg"`
}

type DataResponse[T any] struct {
	BaseResponse
	Data T `json:"data"`
}

type HttpErrorResponse struct {
	Status     string // e.g. "200 OK"
	StatusCode int    `json:"status_code"` // http 状态码
	Body       []byte `json:"body"`        // http 响应体
}

func (h *HttpErrorResponse) Error() string {
	return fmt.Sprintf("response error: statusCode: %d, status: %s", h.StatusCode, h.Status)
}

func (h *HttpErrorResponse) String() string {
	return fmt.Sprintf("HttpErrorResponse: statusCode: %d, status: %s, body: %s", h.StatusCode, h.Status, string(h.Body))
}
