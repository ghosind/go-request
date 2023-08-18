# Golang HTTP请求工具库

![test](https://github.com/ghosind/go-request/workflows/test/badge.svg)
[![codecov](https://codecov.io/gh/ghosind/go-request/branch/main/graph/badge.svg)](https://codecov.io/gh/ghosind/go-request)
![Version Badge](https://img.shields.io/github/v/release/ghosind/go-request)
![License Badge](https://img.shields.io/github/license/ghosind/go-request)
[![Go Reference](https://pkg.go.dev/badge/github.com/ghosind/go-request.svg)](https://pkg.go.dev/github.com/ghosind/go-request)

简体中文 | [English](./README.md)

一个易用的Golang HTTP请求工具库。

- [特性](#特性)
- [安装](#安装)
- [入门](#入门)
  - [`POST`或其他请求](#post-或其他请求)
  - [超时设定](#超时设定)
  - [响应内容处理](#响应内容处理)
- [请求客户端实例](#请求客户端实例)
  - [请求客户端配置](#请求客户端配置)
- [请求配置](#请求配置)

## 特性

- 超时设定或自定义上下文
- 无需手动序列化请求内容，发出请求时自动根据所需格式处理
- 响应内容反序列化封装
- 自动根据响应头部对内容进行解码（解压）

## 安装

> 需要Go 1.18及以上版本支持

可以通过下列命令进行安装：

```sh
go get -u github.com/ghosind/go-request
```

## 入门

下面是一个简单的例子，用于发出一个`GET`请求：

```go
resp, err := request.Request("https://example.com/products/1")
if err != nil {
  // 错误处理
}
// 响应处理
```

### `POST`及其他请求

可以在发出请求时传入请求配置，并设置`Method`属性的值为`"POST"`，以执行发出一个`POST`请求：

```go
resp, err := request.Request("https://example.com/products", RequestConfig{
  Method: "POST",
  Body:   map[string]any{
    "title": "Apple",
  },
})
// 处理错误及响应
```

在没有设置`ContentType`属性的情况下，默认将请求内容序列化为JSON字符串，并将请求头部`Content-Type`的值设置为`application/json`。

除了设置`Method`属性外，还可以直接使用内置的`POST`方法执行`POST`请求，其等同于调用`Request`方法并设置相应的配置：

```go
resp, err := request.POST("https://example.com/products", RequestConfig{
  Body: map[string]any{
    "title": "Apple",
  },
})
// 处理错误及响应
```

除了上述的`POST`方法外，还提供了以下方法，用于执行不同请求方法的请求：

- `DELETE`
- `GET`
- `HEAD`
- `OPTIONS`
- `PATCH`
- `POST`
- `PUT`

> 在使用上述方法而非`Request`，且设置了配置中的`Method`属性的情况下，将会覆盖配置中的设置。

### 超时设定

在默认情况下，所有的请求都将设置一个1秒钟的默认超时时间。若要修改超时时间，可以通过配置中的`Timeout`属性进行调整，其值为以毫秒为单位的整数。

```go
resp, err := request.Request("https://example.com", request.RequestConfig{
  Timeout: 3000, // 3000毫秒，即3秒
})
// 错误及响应处理
```

另外，也可将`Timeout`属性的值设置为`request.RequestTimeoutNone`，用于禁用超时设定。

> 在通过`Context`属性传入自定义上下文的情况下，将不再执行超时的设定。若需要对请求超时进行控制，则需要进行手动处理。

### 响应内容处理

对于响应的内容，提供了`ToObject`及`ToString`方法用于简便处理。例如可以使用`ToString`方法读取响应内容，并以字符串的形式返回。

```go
content, resp, err := ToString(request.Request("https://example.com/products/1"))
if err != nil {
  // 错误处理
}
// content: {"id":1,"title":"iPhone9",...
// 响应处理
```

`ToObject`方法使用时需要设置解析的结构类型，它将读取响应的内容类型，并尝试将响应内容解析至指定类型。

```go
type Product struct {
  ID    int    `json:"id"`
  Title string `json:"title"`
}

product, resp, err := ToObject[Product](request.Request("https://example.com/products/1"))
if err != nil {
  // 错误处理
}
// product: {1 iPhone9}
// 响应处理
```

> `ToObject`与`ToString`方法在执行后都将调用响应体的`Body.Close()`方法。

## 请求客户端实例

对于需要使用一些公用配置（例如相同的请求目标网站、相同的头部值等），可以创建一个请求客户端实例，并传入自定义的配置。例如下面的例子中，将创建一个请求客户端实例并将其基础URL设置为`"https://example.com/"`，随后使用该客户端实例进行请求操作时，都将默认使用该基础URL。

```go
cli := request.New(request.Config{
  BaseURL: "https://example.com/",
})

resp, err := cli.GET("/products/1")
// 错误及响应处理
```

### 请求客户端配置

| 属性 | 类型 | 描述 |
|:-----:|:----:|-------------|
| `BaseURL` | `string` | 基础URL，在请求时将会对其与请求的`url`参数进行拼接，成为最终请求的目标地址。 |
| `Headers` | `map[string][]string` | 自定义头部 |
| `Timeout` | `int` | 以毫秒为单位的超时时长设定 |
| `UserAgent` | `string` | 自定义UserAgent |
| `MaxRedirects` | `int` | 最大跳转次数 |
| `ValidateStatus` | `func(int) bool` | 响应有效性判断方法 |

## 请求配置

下面是可以使用的请求配置，其中所有属性都为可选属性。

| 属性 | 类型 | 描述 |
|:-----:|:----:|-------------|
| `BaseURL` | `string` | 基础URL，在请求时将会对其与请求的`url`参数进行拼接，成为最终请求的目标地址。 |
| `Parameters` | `map[string][]string` | 自定义参数 |
| `Headers` | `map[string][]string` | 自定义请求头部 |
| `Timeout` | `int` | 以毫秒为单位的超时时长设定 |
| `Context` | `context.Context` | 用于请求的上下文 |
| `Body` | `any` | 请求内容 |
| `Method` | `string` | 请求方式，默认为`GET` |
| `ContentType` | `string` | 请求内容类型，当前可用值包括：`"json"`，默认为`"json"` |
| `UserAgent` | `string` | 自定义UserAgent |
| `Auth` | `*BasicAuthConfig` | HTTP Basic Auth设置 |
| `MaxRedirects` | `int` | 最大跳转次数 |
| `ValidateStatus` | `func(int) bool` | 响应有效性判断方法 |
