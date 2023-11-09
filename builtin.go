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

// UseRequestInterceptor adds the request interceptors to the default client. It'll return their ID
// and you can remove these interceptors with the ID by the RemoveRequestInterceptor method.
//
//	UseRequestInterceptor(func (req *http.Request) error {
//		// do something
//		return nil
//	})
func UseRequestInterceptor(interceptors ...RequestInterceptor) []uint64 {
	return defaultClient.UseRequestInterceptor(interceptors...)
}

// RemoveRequestInterceptor removes the request interceptor by the specified interceptor ID, and it
// returns a boolean value to indicate the result.
//
//	ids := UseRequestInterceptor(func (req *http.Request) error {
//		// do something
//		return nil
//	})
//
//	RemoveRequestInterceptor(ids[0])
func RemoveRequestInterceptor(interceptorId uint64) bool {
	return defaultClient.RemoveRequestInterceptor(interceptorId)
}

// UseResponseInterceptor adds the response interceptors to the default client. It'll return their
// ID and you can remove these interceptors with the ID by the RemoveResponseInterceptor method.
//
//	UseResponseInterceptor(func (resp *http.Response) error {
//		// do something
//		return nil
//	})
func UseResponseInterceptor(interceptors ...ResponseInterceptor) []uint64 {
	return defaultClient.UseResponseInterceptor(interceptors...)
}

// RemoveResponseInterceptor removes the response interceptor by the specified interceptor ID, and
// it returns a boolean value to indicate the result.
//
//	ids := UseResponseInterceptor(func (resp *http.Response) error {
//		// do something
//		return nil
//	})
//
//	RemoveResponseInterceptor(ids[0])
func RemoveResponseInterceptor(interceptorId uint64) bool {
	return defaultClient.RemoveResponseInterceptor(interceptorId)
}
