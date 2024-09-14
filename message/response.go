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

import "github.com/chenmingyong0423/go-coze/common/response"

type ListResponse[T any] struct {
	response.DataResponse[T]
	FirstId string `json:"first_id"`
	LastId  string `json:"last_id"`
	HasMore bool   `json:"has_more"`
}

type ModifyResponse[T any] struct {
	response.BaseResponse
	Message T `json:"message"`
}
