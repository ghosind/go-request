package request

import (
	"context"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"
)

type RequestOptions struct {
	// BaseURL will prepended to the url of the request unless the url is absolute.
	BaseURL string
	// Timeout specifies the number of milliseconds before the request times out. This value will be
	// ignored if the `Content` field in the request options is set. It indicates no time-out
	// limitation if the value is -1.
	Timeout int
	// Context is a `context.Content` object that is used for manipulating the request by yourself.
	// The `Timeout` field will be ignored if this value is not empty, and you need to control
	// timeout by yourself.
	Context context.Context
	// Parameters are the URL parameters to be sent with the request.
	Parameters map[string][]string
	// Headers are custom headers to be sent.
	//
	//	resp, err := request.Request("http://example.com", RequestOptions{
	//	  Headers: map[string][]string{
	//	    "Authorization": {"Bearer XXXXX"},
	//	  },
	//	})
	Headers map[string][]string
	// Body is the data to be sent as the request body. It'll be encoded with the content type
	// specified by the `ContentType` field in the request options, or encoded as a JSON if the
	// `ContentType` field is empty. It'll skip the encode processing if the value is a string or a
	// slice of bytes.
	Body any
	// Method indicates the HTTP method of the request, default GET.
	Method string
	// ContentType indicates the type of data that will encode and send to the server. Available
	// options are: "json", default "json".
	ContentType string
	// UserAgent sets the client's User-Agent field in the request header. It'll overwrite the value
	// of the `User-Agent` field in the request headers.
	UserAgent string
	// Auth indicates that HTTP Basic auth should be used. It will set an `Authorization` header,
	// and it'll also overwriting any existing `Authorization` field in the request header.
	//
	//	resp, err := request.Request("https://example.com", RequestOptions{
	//	  Auth: &request.BasicAuthConfig{
	//	    Username: "user",
	//	    Password: "pass",
	//	  },
	//	})
	Auth *BasicAuthConfig
}

// BasicAuthConfig indicates the config of the HTTP Basic Auth that is used for the request.
type BasicAuthConfig struct {
	// Username indicates the username used for HTTP Basic Auth
	Username string
	// Password indicates the password used for HTTP Basic Auth
	Password string
}

// urlPattern is the regular expression pattern for checking whether an URL is starting with HTTP
// or HTTPS protocol or not.
var urlPattern *regexp.Regexp = regexp.MustCompile(`^https?://.+`)

// request creates an HTTP request with the specific HTTP method, the request options, and the
// client config, and send it to the specific destination by the URL.
func (cli *Client) request(method, url string, opts ...RequestOptions) (*http.Response, error) {
	var opt RequestOptions

	if len(opts) > 0 {
		opt = opts[0]
	} else {
		opt = RequestOptions{}
	}

	if method == "" {
		method = opt.Method
	}

	req, canFunc, err := cli.makeRequest(method, url, opt)
	if err != nil {
		return nil, err
	}
	if canFunc != nil {
		defer canFunc()
	}

	httpClient := cli.getHTTPClient()
	defer func() {
		cli.clientPool.Put(httpClient)
	}()

	return httpClient.Do(req)
}

// makeRequest creates a new `http.Request` object with the specific HTTP method, request url, and
// other configurations.
func (cli *Client) makeRequest(method, url string, opt RequestOptions) (*http.Request, context.CancelFunc, error) {
	method, err := cli.getRequestMethod(method)
	if err != nil {
		return nil, nil, err
	}

	url, err = cli.parseURL(url, opt)
	if err != nil {
		return nil, nil, err
	}

	body, err := cli.getRequestBody(opt)
	if err != nil {
		return nil, nil, err
	}

	ctx, canFunc := cli.getContext(opt)

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		canFunc()
		return nil, nil, err
	}

	if err := cli.attachRequestHeaders(req, opt); err != nil {
		canFunc()
		return nil, nil, err
	}

	return req, canFunc, nil
}

// getRequestMethod validates and returns the HTTP method of the request. It'll return "GET" if the
// value of the method is empty.
func (cli *Client) getRequestMethod(method string) (string, error) {
	if method == "" {
		return http.MethodGet, nil
	}

	method = strings.ToUpper(method)
	switch method {
	case http.MethodConnect, http.MethodDelete, http.MethodGet, http.MethodHead, http.MethodOptions,
		http.MethodPatch, http.MethodPost, http.MethodPut, http.MethodTrace:
		return method, nil
	default:
		return "", ErrInvalidMethod
	}
}

// attachRequestHeaders set the field values of the request headers by the request options or
// client configurations. It'll overwrite `Content-Type`, `User-Agent`, and other fields in the
// request headers by the config.
func (cli *Client) attachRequestHeaders(req *http.Request, opt RequestOptions) error {
	cli.setHeaders(req, opt)

	if err := cli.setContentType(req, opt); err != nil {
		return err
	}

	cli.setUserAgent(req, opt)

	if opt.Auth != nil {
		req.SetBasicAuth(opt.Auth.Username, opt.Auth.Password)
	}

	return nil
}

// setHeaders set the field values of the request headers from the request options or the client
// configurations. The fields in the request options will overwrite the same fields in the client
// configuration.
func (cli *Client) setHeaders(req *http.Request, opt RequestOptions) {
	if opt.Headers != nil {
		for k, v := range opt.Headers {
			for _, val := range v {
				req.Header.Add(k, val)
			}
		}
	}

	if cli.Headers != nil {
		for k, v := range cli.Headers {
			if _, existed := req.Header[k]; existed {
				continue
			}

			for _, val := range v {
				req.Header.Add(k, val)
			}
		}
	}
}

// setContentType checks the "Content-Type" field in the request headers, and set it by the
// "ContentType" field value from the request options if no value is set in the headers.
func (cli *Client) setContentType(req *http.Request, opt RequestOptions) error {
	contentType := req.Header.Get("Content-Type")
	if contentType != "" {
		return nil
	}

	switch strings.ToLower(opt.ContentType) {
	case RequestContentTypeJSON, "":
		contentType = "application/json"
	default:
		return ErrUnsupportedType
	}

	req.Header.Set("Content-Type", contentType)

	return nil
}

// setUserAgent checks the user agent value in the request options or the client configurations,
// and set it as the value of the `User-Agent` field in the request headers.
// Default "go-request/x.x".
func (cli *Client) setUserAgent(req *http.Request, opt RequestOptions) {
	userAgent := opt.UserAgent
	if userAgent == "" && cli.UserAgent != "" {
		userAgent = cli.UserAgent
	}

	if userAgent == "" {
		userAgent = RequestDefaultUserAgent
	}

	req.Header.Set("User-Agent", userAgent)
}

func (cli *Client) parseURL(uri string, opt RequestOptions) (string, error) {
	baseURL, extraPath, err := cli.getURL(uri, opt)
	if err != nil {
		return "", err
	}

	obj, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	if extraPath != "" {
		obj.Path = path.Join(obj.Path, extraPath)
	}

	if opt.Parameters != nil {
		// attach parameters to request url
		query := obj.Query()

		for k, vv := range opt.Parameters {
			for _, v := range vv {
				query.Add(k, v)
			}
		}

		obj.RawQuery = query.Encode()
	}

	return obj.String(), nil
}

// getURL returns the base url and extra path components from url parameter, optional config, and
// instance config.
func (cli *Client) getURL(url string, opt RequestOptions) (string, string, error) {
	if url != "" && urlPattern.MatchString(url) {
		return url, "", nil
	}

	baseURL := opt.BaseURL
	if baseURL == "" && cli.BaseURL != "" {
		baseURL = cli.BaseURL
	}
	if baseURL == "" {
		baseURL = url
		url = ""
	}

	if baseURL == "" {
		return "", "", ErrNoURL
	}

	if !urlPattern.MatchString(baseURL) {
		// prepend https as scheme if no scheme part in the url.
		baseURL = "https://" + baseURL
	}

	return baseURL, url, nil
}

// getContext creates a Context by the request options or client settings, or returns the Context
// that is set in the request options.
func (cli *Client) getContext(opt RequestOptions) (context.Context, context.CancelFunc) {
	if opt.Context != nil {
		return opt.Context, nil
	}

	baseCtx := context.Background()

	timeout := RequestTimeoutDefault
	if opt.Timeout > 0 || opt.Timeout == RequestTimeoutNoLimit {
		timeout = opt.Timeout
	} else if cli.timeout > 0 || cli.timeout == RequestTimeoutNoLimit {
		timeout = cli.timeout
	}

	if timeout == RequestTimeoutNoLimit {
		return baseCtx, nil
	} else {
		return context.WithTimeout(baseCtx, time.Duration(timeout)*time.Millisecond)
	}
}
