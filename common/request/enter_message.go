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

package request

type EnterMessage struct {
	// The role who returns the message.
	// 发送这条消息的实体。
	Role string `json:"role"`
	// The type of the message when the role is assistant.
	// 当 role = assistant 时，用于标识 Bot 的消息类型。
	Type string `json:"type,omitempty"`
	// The returned content.
	// 消息内容。
	Content string `json:"content,omitempty"`
	// The type of the return content.
	// 消息内容的类型。
	ContentType string `json:"content_type,omitempty"`

	// Additional information when creating a message, and this additional information will also be returned when retrieving messages.
	// 创建消息时的附加消息，获取消息时也会返回此附加消息。
	MetaData map[string]any `json:"meta_data,omitempty"`
}

type EnterMessageBuilder struct {
	enterMessage EnterMessage
}

func NewEnterMessageBuilder() *EnterMessageBuilder {
	return &EnterMessageBuilder{}
}

func (b *EnterMessageBuilder) Role(role string) *EnterMessageBuilder {
	b.enterMessage.Role = role
	return b
}

func (b *EnterMessageBuilder) Type(type_ string) *EnterMessageBuilder {
	b.enterMessage.Type = type_
	return b
}

func (b *EnterMessageBuilder) Content(content string) *EnterMessageBuilder {
	b.enterMessage.Content = content
	return b
}

func (b *EnterMessageBuilder) ContentType(contentType string) *EnterMessageBuilder {
	b.enterMessage.ContentType = contentType
	return b
}

func (b *EnterMessageBuilder) MetaData(metaData map[string]any) *EnterMessageBuilder {
	b.enterMessage.MetaData = metaData
	return b
}

func (b *EnterMessageBuilder) Build() EnterMessage {
	return b.enterMessage
}
