package request

import "net/http"

var defaultClient *Client

func init() {
	defaultClient = New()
}

func Request(url string, opts ...RequestOptions) (*http.Response, error) {
	return defaultClient.Request(url, opts...)
}

func DELETE(url string, opt ...RequestOptions) (*http.Response, error) {
	return defaultClient.DELETE(url, opt...)
}

func GET(url string, opt ...RequestOptions) (*http.Response, error) {
	return defaultClient.GET(url, opt...)
}

func HEAD(url string, opt ...RequestOptions) (*http.Response, error) {
	return defaultClient.HEAD(url, opt...)
}

func OPTIONS(url string, opt ...RequestOptions) (*http.Response, error) {
	return defaultClient.OPTIONS(url, opt...)
}

func PATCH(url string, opt ...RequestOptions) (*http.Response, error) {
	return defaultClient.PATCH(url, opt...)
}

func POST(url string, opt ...RequestOptions) (*http.Response, error) {
	return defaultClient.POST(url, opt...)
}

func PUT(url string, opt ...RequestOptions) (*http.Response, error) {
	return defaultClient.PUT(url, opt...)
}
