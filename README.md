# 

<h1 align="center">
  go-coze
</h1>

[![GitHub Repo stars](https://img.shields.io/github/stars/chenmingyong0423/go-coze)](https://github.com/chenmingyong0423/go-coze/stargazers)
[![GitHub issues](https://img.shields.io/github/issues/chenmingyong0423/go-coze)](https://github.com/chenmingyong0423/go-coze/issues)
[![GitHub License](https://img.shields.io/github/license/chenmingyong0423/go-coze)](https://github.com/chenmingyong0423/go-coze/blob/main/LICENSE)
[![GitHub release (with filter)](https://img.shields.io/github/v/release/chenmingyong0423/go-coze)](https://github.com/chenmingyong0423/go-coze)
[![Go Report Card](https://goreportcard.com/badge/github.com/chenmingyong0423/go-coze)](https://goreportcard.com/report/github.com/chenmingyong0423/go-coze)
[![All Contributors](https://img.shields.io/badge/all_contributors-1-orange.svg?style=flat-square)](#contributors-)

`go-coze` 库是一个用于简化 `Coze API` 调用的库，通过这个库，开发者可以更高效，更简洁地与 `Coze API` 交互，减少重复代码，提高开发效率。

## 功能
- **链式调用**
    - 通过链式调用的方式封装请求参数和调用接口，使代码更加简洁和可读。

- **非流式 API 交互**
    - 适用于一次性获取数据的场景。例如当 `stream` 参数指定为 `false` 的场景。

- **流式 API 交互**
    - 支持处理流式响应，例如当 `stream` 参数被指定为 `true` 的场景。

# 入门指南
## 安装
```shell
go get github.com/chenmingyong0423/go-coze
```

## 使用
### 非流式 API 交互
```go
// 创建一个聊天对象
chat := NewChat("botID", "user", "personalAccessToken")

// 创建新的会话对象并设置会话流和类型
session := chat.Chat(false)

// 添加请求参数并发送以及处理错误
resp, err := session.WithQuery("你好").
  WithConversationId("conversationId").
  WithChatHistory(nil).
  WithCustomVariables(nil).
  Request(context.Background())
```
非流式 `API` 交互需要调用 `Request` 方法，该方法会返回一个 `NonStreamingResponse` 对象和一个 `error` 对象。
### 流式 API 交互
```go
// 创建一个聊天对象
chat := NewChat("botID", "user", "personalAccessToken")

// 创建新的会话对象并设置会话流和类型
session := chat.Chat(true)

// 添加请求参数并发送以及处理错误
respChan, errChan := session.WithQuery("你好").
  WithConversationId("conversationId").
  WithChatHistory(nil).
  WithCustomVariables(nil).
  Request(context.Background())
for {
    select {
    case resp, ok := <-respChan:
        if !ok {
            respChan = nil
        } else {
            fmt.Println(resp)
        }
    case err, ok := <-errChan:
        if !ok {
            errChan = nil
        } else {
            panic(err)
        }
    }
    if respChan == nil && errChan == nil {
        break
    }
}
```
流式 `API` 交互需要调用 `StreamRequest` 方法，该方法会返回一个 `chan *StreamingResponse` 对象和一个 `chan error` 对象。