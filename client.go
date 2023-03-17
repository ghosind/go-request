package request

import (
	"net/http"
	"sync"
)

type Client struct {
	// BaseURL will be prepended to all request URL unless URL is absolute.
	BaseURL string
	// Headers are custom headers to be sent.
	Headers http.Header

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
}

const (
	// RequestTimeoutDefault is the default timeout for request.
	RequestTimeoutDefault int = 1000
	// RequestTimeoutNone means
	RequestTimeoutNone int = -1
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
		cli.setHeader(cfg.Headers)
	}

	return cli
}

// setHeader initializes client's Headers field from config.
func (cli *Client) setHeader(headers map[string][]string) {
	for k, v := range headers {
		if len(v) > 0 {
			cli.Headers[k] = make([]string, len(v))
			copy(cli.Headers[k], v)
		} else {
			cli.Headers[k] = nil
		}
	}
}

// getHTTPClient returns a http.Client from the pool.
func (cli *Client) getHTTPClient() *http.Client {
	return cli.clientPool.Get().(*http.Client)
}
