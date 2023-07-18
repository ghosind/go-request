package request

import (
	"context"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"time"
)

type RequestOptions struct {
	BaseURL     string
	Timeout     int
	Context     context.Context
	Parameters  map[string][]string
	Headers     map[string][]string
	Body        any
	Method      string
	ContentType string
	// UserAgent sets the client's User-Agent field in the request header.
	UserAgent string
}

var urlPattern *regexp.Regexp = regexp.MustCompile(`^https?://`)

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

func (cli *Client) makeRequest(method, url string, opt RequestOptions) (*http.Request, context.CancelFunc, error) {
	if method == "" {
		method = http.MethodGet
	}

	url, err := cli.parseURL(url, opt)
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

	cli.addHeadersToRequest(req, opt)

	return req, canFunc, nil
}

func (cli *Client) addHeadersToRequest(req *http.Request, opt RequestOptions) {
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

	if opt.ContentType != "" {
		req.Header.Set("Content-Type", opt.ContentType)
	}
	if ct := req.Header.Get("Content-Type"); ct == "" {
		req.Header.Set("Content-Type", "application/json") // Set default content type
	}

	userAgent := opt.UserAgent
	if userAgent == "" && cli.UserAgent != "" {
		userAgent = cli.UserAgent
	}
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}
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

func (cli *Client) getContext(opt RequestOptions) (context.Context, context.CancelFunc) {
	if opt.Context != nil {
		return opt.Context, nil
	}

	baseCtx := context.Background()

	timeout := RequestTimeoutDefault
	if opt.Timeout > 0 || opt.Timeout == RequestTimeoutNone {
		timeout = opt.Timeout
	} else if cli.timeout > 0 || cli.timeout == RequestTimeoutNone {
		timeout = cli.timeout
	}

	if timeout == RequestTimeoutNone {
		return baseCtx, nil
	} else {
		return context.WithTimeout(baseCtx, time.Duration(timeout)*time.Millisecond)
	}
}
