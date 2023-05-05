package request

import "net/http"

var defaultClient *Client

func init() {
	defaultClient = New()
}

// Request the specific url with the optional request settings by the default
// client instance, and decode the response to the out value.
func Request(url string, opts ...RequestOptions) (*http.Response, error) {
	return defaultClient.Request(url, opts...)
}
