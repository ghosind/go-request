package request

import (
	"net"
	"net/http"
	"net/url"
	"sync"
)

// Client is the HTTP requesting client.
type Client struct {
	// BaseURL will be prepended to all request URL unless URL is absolute.
	BaseURL string
	// Headers are custom headers to be sent.
	Headers map[string][]string
	// MaxRedirects defines the maximum number of redirects for this client, default 5.
	MaxRedirects int
	// Parameters are the parameters to be sent.
	Parameters map[string][]string
	// ParametersSerializer is a function to charge of serializing the URL query parameters.
	ParametersSerializer func(map[string][]string) string
	// Proxy is the config of the proxy server.
	Proxy *ProxyConfig
	// Timeout specifies the time before the request times out.
	Timeout int
	// UserAgent sets the client's User-Agent field in the request header.
	UserAgent string
	// ValidateStatus defines whether the status code of the response is valid or not, and it'll
	// return an error if fails to validate the status code.
	ValidateStatus func(int) bool

	// clientPool is for save http.Client instances.
	clientPool *sync.Pool
}

// Config is the config for the HTTP requesting client.
type Config struct {
	// BaseURL will be prepended to all request URL unless URL is absolute.
	BaseURL string
	// Headers are custom headers to be sent, and they'll be overwritten if the
	// same key is presented in the request.
	Headers map[string][]string
	// MaxRedirects defines the maximum number of redirects for this client, default 5.
	MaxRedirects int
	// Parameters are the parameters to be sent for all requests of the client. It will be
	// overwritten if the same key is in the parameters of the request options.
	Parameters map[string][]string
	// ParametersSerializer is a function to charge of serializing the URL query parameters.
	ParametersSerializer func(map[string][]string) string
	// Proxy defines the address and the auth credentials of the proxy server, it will be overwritten
	// by the request options' proxy config if the proxy config is not empty in the request options.
	//
	// You can also define the proxy by the `http_proxy` and `https_proxy` environment variables. If
	// no proxy config in the request options or the client config, the request will try to get a
	// proxy from the environment variables.
	Proxy *ProxyConfig
	// Timeout is request timeout in milliseconds.
	Timeout int
	// UserAgent sets the client's User-Agent field in the request header.
	UserAgent string
	// ValidateStatus defines whether the status code of the response is valid or not, and it'll
	// return an error if fails to validate the status code. Default, it sets the result to fail if
	// the status code is less than 200, or greater than and equal to 400.
	//
	//	cli := request.New(request.Config{
	//	  ValidateStatus: func (status int) bool {
	//	    // Only success if the status code of response is 2XX
	//	    return status >= http.StatusOk && status <= http.StatusMultipleChoices
	//	  },
	//	})
	ValidateStatus func(int) bool
}

const (
	// RequestTimeoutDefault is the default timeout for request.
	RequestTimeoutDefault int = 1000
	// RequestTimeoutNoLimit means no timeout limitation.
	RequestTimeoutNoLimit int = -1

	// RequestMaxRedirects is the default maximum number of redirects.
	RequestDefaultMaxRedirects int = 5
	// RequestNoRedirect means it'll never redirect automatically.
	RequestNoRedirects int = -1

	// RequestDefaultUserAgent is the default user agent for all requests that are sent by this
	// package.
	RequestDefaultUserAgent string = "go-request/0.2"
)

// New creates and returns a new Client instance.
//
//	cli := request.New(request.Config{
//	  BaseURL: "https://www.example.com", // the base URL for every request sent by this client
//	  Timeout: 5000, // the default timeout for this client
//	  // ...
//	})
func New(config ...Config) *Client {
	cli := new(Client)

	cli.clientPool = &sync.Pool{
		New: func() any {
			return new(http.Client)
		},
	}

	cli.Headers = make(http.Header)
	cli.Parameters = make(url.Values)

	if len(config) > 0 {
		cfg := config[0]

		cli.BaseURL = cfg.BaseURL
		cli.MaxRedirects = cfg.MaxRedirects
		cli.ParametersSerializer = cfg.ParametersSerializer
		cli.Proxy = cfg.Proxy
		cli.Timeout = cfg.Timeout
		cli.UserAgent = cfg.UserAgent
		cli.ValidateStatus = cfg.ValidateStatus

		cli.initClientHeaders(cfg.Headers)
		cli.initClientParameters(cfg.Parameters)
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

// Req creates a request with chaining API, and sets the destination to `url`.
//
//	resp, err := cli.Req("http://example.com").
//	  POST().
//	  Body(data).
//	  SetHeader("Accept-Encoding", "gzip").
//	  Do()
func (cli *Client) Req(url string) *RequestOptions {
	opt := new(RequestOptions)

	opt.url = url
	opt.client = cli

	return opt
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

// initClientParameters initializes client's Parameters field from config.
func (cli *Client) initClientParameters(parameters map[string][]string) {
	for k, v := range parameters {
		if len(v) > 0 {
			cli.Parameters[k] = make([]string, len(v))
			copy(cli.Parameters[k], v)
		}
	}
}

// getHTTPClient gets an `http.Client` from the pool, and resets it to default state.
func (cli *Client) getHTTPClient(opt RequestOptions) *http.Client {
	if cli.clientPool == nil {
		cli.clientPool = &sync.Pool{
			New: func() any {
				return new(http.Client)
			},
		}
	}

	httpClient := cli.clientPool.Get().(*http.Client)

	maxRedirects := opt.MaxRedirects
	if maxRedirects == 0 {
		maxRedirects = cli.MaxRedirects
	}
	if maxRedirects < RequestNoRedirects || maxRedirects == 0 {
		maxRedirects = RequestDefaultMaxRedirects
	}

	httpClient.CheckRedirect = cli.getCheckRedirect(maxRedirects)
	httpClient.Transport = cli.getTransport(opt)

	return httpClient
}

// getCheckRedirect returns a new check redirects handler for `http.Client`. This function will
// never return errors except `http.ErrUseLastResponse` error that the redirects number is greater
// than the maximum limitation.
func (cli *Client) getCheckRedirect(
	maxRedirects int,
) func(req *http.Request, via []*http.Request) error {
	if maxRedirects == RequestDefaultMaxRedirects {
		return cli.defaultCheckRedirect
	}

	return func(req *http.Request, via []*http.Request) error {
		if len(via) >= maxRedirects {
			return http.ErrUseLastResponse
		}

		return nil
	}
}

// defaultCheckRedirect is the default redirect check handler, and it returns
// `http.ErrUseLastResponse` if the number of redirects is greater than or equal to the default
// number of maximum redirects. It returns `http.ErrUseLastResponse` to terminate the redirection
// but does not return an error.
func (cli *Client) defaultCheckRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= RequestDefaultMaxRedirects {
		return http.ErrUseLastResponse
	}

	return nil
}

// defaultValidateStatus is the default handler to check the status code of the responses. It only
// returns true if the status code is greater than or equal to 200, and less than 400.
func (cli *Client) defaultValidateStatus(status int) bool {
	return status >= http.StatusOK && status < http.StatusBadRequest
}

// getTransport gets the transport by the request options.
func (cli *Client) getTransport(opt RequestOptions) http.RoundTripper {
	proxy := cli.getProxy(opt)
	if proxy == nil {
		return nil
	}

	return &http.Transport{
		Proxy: proxy,
	}
}

// getProxy tries to get a proxy by the proxy config from the request options, and it returns
// `http.ProxyFromEnvironment` if there is no proxy config in the request options.
func (cli *Client) getProxy(opt RequestOptions) func(*http.Request) (*url.URL, error) {
	proxy := opt.Proxy
	if proxy == nil {
		proxy = cli.Proxy
	}
	if proxy == nil {
		return http.ProxyFromEnvironment
	}

	proxyUrl := new(url.URL)
	proxyUrl.Scheme = proxy.Protocol
	proxyUrl.Host = net.JoinHostPort(proxy.Host, proxy.Port)
	if proxy.Username != "" || proxy.Password != "" {
		proxyUrl.User = url.UserPassword(proxy.Username, proxy.Password)
	}

	return http.ProxyURL(proxyUrl)
}
