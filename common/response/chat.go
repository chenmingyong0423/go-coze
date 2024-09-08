// Copyright 2024 chenmingyong0423

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

type Chat struct {
	Id             string            `json:"id"`
	ConversationId string            `json:"conversation_id"`
	BotId          string            `json:"bot_id"`
	CreatedAt      int64             `json:"created_at,omitempty"`
	CompletedAt    int64             `json:"completed_at,omitempty"`
	FailedAt       int64             `json:"failed_at,omitempty"`
	MetaData       map[string]string `json:"meta_data,omitempty"`
	LastError      LastError         `json:"last_error,omitempty"`
	Status         string            `json:"status"`
	RequiredAction RequiredAction    `json:"required_action,omitempty"`
	Usage          Usage             `json:"usage,omitempty"`
}

type LastError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type RequiredAction struct {
	Type              string           `json:"type,omitempty"`
	SubmitToolOutputs SubmitToolOutput `json:"submit_tool_outputs,omitempty"`
}

type SubmitToolOutput struct {
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type ToolCall struct {
	Id       string   `json:"id,omitempty"`
	Type     string   `json:"type,omitempty"`
	Function Function `json:"function,omitempty"`
}

type Function struct {
	Name     string `json:"name,omitempty"`
	Argument string `json:"argument,omitempty"`
}

type Usage struct {
	TokenCount  int `json:"token_count,omitempty"`
	OutputCount int `json:"output_count,omitempty"`
	InputCount  int `json:"input_count,omitempty"`
}
