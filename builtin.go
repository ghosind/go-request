package request

import "net/http"

// defaultClient is the default HTTP client for sending requests without creating a new client
// object.
var defaultClient *Client

func init() {
	defaultClient = New()
}

// Request performs an HTTP request to the specific URL with the request options. If no request
// options are set, it will be sent as an HTTP GET request.
//
//	resp, err := request.Request("https://example.com")
//	if err != nil {
//	  // Error handling
//	}
//	// Response handling
func Request(url string, opts ...RequestOptions) (*http.Response, error) {
	return defaultClient.Request(url, opts...)
}

// DELETE performs an HTTP DELETE request to the specific URL with the request options.
func DELETE(url string, opt ...RequestOptions) (*http.Response, error) {
	return defaultClient.DELETE(url, opt...)
}

// GET performs an HTTP GET request to the specific URL with the request options.
func GET(url string, opt ...RequestOptions) (*http.Response, error) {
	return defaultClient.GET(url, opt...)
}

// HEAD performs an HTTP HEAD request to the specific URL with the request options.
func HEAD(url string, opt ...RequestOptions) (*http.Response, error) {
	return defaultClient.HEAD(url, opt...)
}

// OPTIONS performs an HTTP OPTIONS request to the specific URL with the request options.
func OPTIONS(url string, opt ...RequestOptions) (*http.Response, error) {
	return defaultClient.OPTIONS(url, opt...)
}

// PATCH performs an HTTP PATCH request to the specific URL with the request options.
func PATCH(url string, opt ...RequestOptions) (*http.Response, error) {
	return defaultClient.PATCH(url, opt...)
}

// POST performs an HTTP POST request to the specific URL with the request options.
func POST(url string, opt ...RequestOptions) (*http.Response, error) {
	return defaultClient.POST(url, opt...)
}

// PUT performs an HTTP PUT request to the specific URL with the request options.
func PUT(url string, opt ...RequestOptions) (*http.Response, error) {
	return defaultClient.PUT(url, opt...)
}

// Req creates a request with chaining API, and sets the destination to `url`.
//
//	resp, err := request.Req("http://example.com").
//	  POST().
//	  Body(data).
//	  SetHeader("Accept-Encoding", "gzip").
//	  Do()
func Req(url string) *RequestOptions {
	return defaultClient.Req(url)
}
