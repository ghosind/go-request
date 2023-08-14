package request

import (
	"net/http"
	"sync"
)

type Client struct {
	// BaseURL will be prepended to all request URL unless URL is absolute.
	BaseURL string
	// Headers are custom headers to be sent.
	Headers map[string][]string
	// UserAgent sets the client's User-Agent field in the request header.
	UserAgent string

	// clientPool is for save http.Client instances.
	clientPool sync.Pool
	// timeout specifies the time before the request times out.
	timeout int
}

type Config struct {
	// BaseURL will be prepended to all request URL unless URL is absolute.
	BaseURL string
	// Timeout is request timeout in milliseconds.
	Timeout int
	// Headers are custom headers to be sent, and they'll be overwritten if the
	// same key is presented in the request.
	Headers map[string][]string
	// UserAgent sets the client's User-Agent field in the request header.
	UserAgent string
}

const (
	// RequestTimeoutDefault is the default timeout for request.
	RequestTimeoutDefault int = 1000
	// RequestTimeoutNoLimit means no timeout limitation.
	RequestTimeoutNoLimit int = -1

	// RequestDefaultUserAgent is the default user agent for all requests that are sent by this
	// package.
	RequestDefaultUserAgent string = "go-request/0.2"
)

// New creates and returns a new Client instance.
func New(config ...Config) *Client {
	cli := new(Client)

	cli.Headers = make(http.Header)
	cli.clientPool = sync.Pool{
		New: func() any {
			return new(http.Client)
		},
	}

	if len(config) > 0 {
		cfg := config[0]

		cli.BaseURL = cfg.BaseURL
		cli.timeout = cfg.Timeout
		cli.UserAgent = cfg.UserAgent
		cli.initClientHeaders(cfg.Headers)
	}

	return cli
}

// Request performs an HTTP request to the specific URL with the request options and the client
// config. If no request options are set, it will be sent as an HTTP GET request.
//
//	resp, err := cli.Request("https://example.com")
//	if err != nil {
//	  // Error handling
//	}
//	// Response handling
func (cli *Client) Request(url string, opt ...RequestOptions) (*http.Response, error) {
	return cli.request("", url, opt...)
}

// DELETE performs an HTTP DELETE request to the specific URL with the request options and the
// client config.
func (cli *Client) DELETE(url string, opt ...RequestOptions) (*http.Response, error) {
	return cli.request(http.MethodDelete, url, opt...)
}

// GET performs an HTTP GET request to the specific URL with the request options and the client
// config.
func (cli *Client) GET(url string, opt ...RequestOptions) (*http.Response, error) {
	return cli.request(http.MethodGet, url, opt...)
}

// HEAD performs an HTTP HEAD request to the specific URL with the request options and the client
// config.
func (cli *Client) HEAD(url string, opt ...RequestOptions) (*http.Response, error) {
	return cli.request(http.MethodHead, url, opt...)
}

// OPTIONS performs an HTTP OPTIONS request to the specific URL with the request options and the
// client config.
func (cli *Client) OPTIONS(url string, opt ...RequestOptions) (*http.Response, error) {
	return cli.request(http.MethodOptions, url, opt...)
}

// PATCH performs an HTTP PATCH request to the specific URL with the request options and the
// client config.
func (cli *Client) PATCH(url string, opt ...RequestOptions) (*http.Response, error) {
	return cli.request(http.MethodPatch, url, opt...)
}

// POST performs an HTTP POST request to the specific URL with the request options and the client
// config.
func (cli *Client) POST(url string, opt ...RequestOptions) (*http.Response, error) {
	return cli.request(http.MethodPost, url, opt...)
}

// PUT performs an HTTP PUT request to the specific URL with the request options and the client
// config.
func (cli *Client) PUT(url string, opt ...RequestOptions) (*http.Response, error) {
	return cli.request(http.MethodPut, url, opt...)
}

// initClientHeaders initializes client's Headers field from config.
func (cli *Client) initClientHeaders(headers map[string][]string) {
	for k, v := range headers {
		if len(v) > 0 {
			cli.Headers[k] = make([]string, len(v))
			copy(cli.Headers[k], v)
		}
	}
}

// getHTTPClient returns a http.Client from the pool.
func (cli *Client) getHTTPClient() *http.Client {
	return cli.clientPool.Get().(*http.Client)
}
