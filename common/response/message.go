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

type Message struct {
	Id             string         `json:"id"`
	ConversationId string         `json:"conversation_id"`
	BotId          string         `json:"bot_id"`
	ChatId         string         `json:"chat_id"`
	MetaData       map[string]any `json:"meta_data"`
	Role           string         `json:"role"`
	Content        string         `json:"content"`
	ContentType    string         `json:"content_type"`
	CreateTime     int64          `json:"create_time"`
	UpdateTime     int64          `json:"update_time"`
	Type           string         `json:"type"`
}
