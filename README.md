# HTTP Request tool for Go

![test](https://github.com/ghosind/go-request/workflows/test/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/ghosind/go-request)](https://goreportcard.com/report/github.com/ghosind/go-request)
[![codecov](https://codecov.io/gh/ghosind/go-request/branch/main/graph/badge.svg)](https://codecov.io/gh/ghosind/go-request)
![Version Badge](https://img.shields.io/github/v/release/ghosind/go-request)
![License Badge](https://img.shields.io/github/license/ghosind/go-request)
[![Go Reference](https://pkg.go.dev/badge/github.com/ghosind/go-request.svg)](https://pkg.go.dev/github.com/ghosind/go-request)

English | [简体中文](./README_CN.md)

An easy-to-use HTTP request tool for Golang.

- [Features](#features)
- [Installation](#installation)
- [Getting Started](#getting-started)
  - [`POST` and other requests](#post-and-other-requests)
  - [Timeouts](#timeouts)
  - [Response body handling](#response-body-handling)
- [Client Instance](#client-instance)
  - [Client Instance Config](#client-instance-config)
- [Request Config](#request-config)
- [Roadmap](#roadmap)

## Features

- Timeouts or self-control context.
- Serialize request body automatically.
- Response body deserialization wrapper.
- Decode the compressed response body automatically.
- Chaining API.
- Request and Response interceptors.

## Installation

> This package requires Go 1.18 and later versions.

You can install this package by the following command.

```sh
go get -u github.com/ghosind/go-request
```

## Getting Started

The is a minimal example of performing a `GET` request:

```go
resp, err := request.Request("https://example.com/products/1")
if err != nil {
  // handle error
}
// handle response
```

### `POST` and other requests

You can perform a `POST` request with a request config that set the `Method` field's value to `POST`.

```go
resp, err := request.Request("https://example.com/products", RequestConfig{
  Method: "POST",
  Body:   map[string]any{
    "title": "Apple",
  },
})
// handle error or response
```

If the `ContentType` field in the request config is empty, the body data will serialize to a JSON string default, and it'll also set the `Content-Type` field value in the request headers to `application/json`.

You can also use `POST` method to perform a `POST` request with the specific body data.

```go
resp, err := request.POST("https://example.com/products", RequestConfig{
  Body: map[string]any{
    "title": "Apple",
  },
})
// handle error or response
```

We also provided the following methods for performing HTTP requests:

- `DELETE`
- `GET`
- `HEAD`
- `OPTIONS`
- `PATCH`
- `POST`
- `PUT`

> The above methods will overwrite the `Method` field in the request config.

### Timeouts

All the requests will set timeout to 1-second default, you can set a custom timeout value in milliseconds to a request:

```go
resp, err := request.Request("https://example.com", request.RequestConfig{
  Timeout: 3000, // 3 seconds
})
// handle error or response
```

You can also set `Timeout` to `request.RequestTimeoutNone` to disable the timeout mechanism.

> The timeout will be disabled if you set `Context` in the request config, you need to handle it manually.

### Chaining API

You can also make a request by chaining API: 

```go
resp, err := request.Req("http://example.com").
  POST().
  SetBody(map[string]any{ "title": "Apple" }).
  SetTimeout(3000).
  Do()
```

### Response body handling

We provided `ToObject` and `ToString` methods to handle response body. For example, the `ToString` method will read all data in the response body, and return it that represented in a string value.

```go
content, resp, err := ToString(request.Request("https://example.com/products/1"))
if err != nil {
  // handle error
}
// handle response
```

The `ToObject` method will read the content type of response and deserialize the body to the specific type.

```go
type Product struct {
  ID    int    `json:"id"`
  Title string `json:"title"`
}

product, resp, err := ToObject[Product](request.Request("https://example.com/products/1"))
if err != nil {
  // handle error
}
// handle response
```

> Both `ToObject` and `ToString` methods will close the `Body` of the response after reading all data.

## Client Instance

You can create a new client instance with a custom config.

```go
cli := request.New(request.Config{
  BaseURL: "https://example.com/",
})

resp, err := cli.GET("/products/1")
// handle error or response
```

### Client Instance Config

| Field | Type | Description |
|:-----:|:----:|-------------|
| `BaseURL` | `string` | The base url for all requests that performing by this client instance. |
| `Headers` | `map[string][]string` | Custom headers to be sent. |
| `MaxRedirects` | `int` | The maximum number of redirects for this client, default 5. |
| `Parameters` | `map[string][]string` | Custom query string parameters to be sent. |
| `Timeout` | `int` | Timeout in milliseconds. |
| `UserAgent` | `string` | Custom user agent value. |
| `ValidateStatus` | `func(int) bool` | The function checks whether the status code of the response is valid or not. |

## Request Config

There are the available config options for performing a request, and all fields are optional.

| Field | Type | Description |
|:-----:|:----:|-------------|
| `Auth` | `*BasicAuthConfig` | HTTP Basic Auth config. |
| `BaseURL` | `string` | The base url for this requests. |
| `Body` | `any` | The request body. |
| `ContentType` | `string` | The content type of this request. Available options are: `"json"`, and default `"json"`. |
| `Context` | `context.Context` | Self-control context. |
| `DisableDecompress` | `bool` | Indicates whether or not disable decompression of the response body automatically. |
| `Headers` | `map[string][]string` | Custom headers to be sent. |
| `MaxRedirects` | `int` | The maximum number of redirects for the request, default 5. |
| `Method` | `string` | HTTP request method, default `GET`. |
| `Parameters` | `map[string][]string` | Custom query string parameters to be sent. |
| `Timeout` | `int` | Timeout in milliseconds. |
| `UserAgent` | `string` | Custom user agent value. |
| `ValidateStatus` | `func(int) bool` | The function checks whether the status code of the response is valid or not. |
