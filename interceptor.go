package request

import "net/http"

// RequestInterceptor is a function to intercept the requests.
type RequestInterceptor func(*http.Request) error

// ResponseInterceptor is a function to intercept the responses.
type ResponseInterceptor func(*http.Response) error
