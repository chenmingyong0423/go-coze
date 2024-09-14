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

type ObjectString struct {
	// The content type of the multimodal message.
	// 多模态消息内容类型。
	Type string `json:"type"`

	// Text content. Required when type is text.
	// 文本内容。
	Text string `json:"text,omitempty"`

	// The ID of the file or image content.
	// 文件或图片内容的 ID。
	FileId string `json:"file_id,omitempty"`

	// The online address of the file or image content.<br>Must be a valid address that is publicly accessible.
	// file_id or file_url must be specified when type is file or image.
	// 文件或图片内容的在线地址。必须是可公共访问的有效地址。
	// 在 type 为 file 或 image 时，file_id 和 file_url 应至少指定一个。
	FileUrl string `json:"file_url,omitempty"`
}

type ObjectStringBuilder struct {
	objectString ObjectString
}

func NewObjectStringBuilder() *ObjectStringBuilder {
	return &ObjectStringBuilder{}
}

func (b *ObjectStringBuilder) Type(typ string) *ObjectStringBuilder {
	b.objectString.Type = typ
	return b
}

func (b *ObjectStringBuilder) Text(text string) *ObjectStringBuilder {
	b.objectString.Text = text
	return b
}

func (b *ObjectStringBuilder) FileId(fileId string) *ObjectStringBuilder {
	b.objectString.FileId = fileId
	return b
}

func (b *ObjectStringBuilder) FileUrl(fileUrl string) *ObjectStringBuilder {
	b.objectString.FileUrl = fileUrl
	return b
}

func (b *ObjectStringBuilder) Build() ObjectString {
	return b.objectString
}
