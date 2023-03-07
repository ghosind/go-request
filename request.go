package request

import (
	"context"
	"net/http"
	"net/url"
	"path"
	"time"
)

type RequestOptions struct {
	BaseURL    string
	Method     string
	Timeout    int
	Context    context.Context
	Parameters map[string][]string
	Headers    map[string][]string
}

func Request(url string, out any, opts ...RequestOptions) error {
	return defaultClient.Request(url, out, opts...)
}

func (cli *Client) Request(url string, out any, opts ...RequestOptions) error {
	_, err := cli.RequestRaw(url, opts...)
	if err != nil {
		return err
	}

	// TODO: decode body

	return nil
}

func (cli *Client) RequestRaw(url string, opts ...RequestOptions) (*http.Response, error) {
	var opt RequestOptions

	if len(opts) > 0 {
		opt = opts[0]
	} else {
		opt = RequestOptions{}
	}

	req, canFunc, err := cli.makeRequest(url, opt)
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

func (cli *Client) makeRequest(url string, opt RequestOptions) (*http.Request, context.CancelFunc, error) {
	method := "GET"
	if opt.Method != "" {
		method = opt.Method
	}

	url, err := cli.parseURL(url, opt)
	if err != nil {
		return nil, nil, err
	}

	ctx, canFunc := cli.getContext(opt)

	// TODO: body

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
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
			for _, val := range v {
				req.Header.Add(k, val)
			}
		}
	}
}

func (cli *Client) parseURL(uri string, opt RequestOptions) (string, error) {
	baseURL, uri, err := cli.getURL(uri, opt)
	if err != nil {
		return "", err
	}

	obj, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	if uri != "" {
		obj.Path = path.Join(obj.Path, uri)
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

func (cli *Client) getURL(uri string, opt RequestOptions) (string, string, error) {
	baseURL := opt.BaseURL
	if baseURL == "" && cli.BaseURL != "" {
		baseURL = cli.BaseURL
	}
	if baseURL == "" {
		baseURL = uri
		uri = ""
	}

	if baseURL == "" {
		return "", "", ErrNoURL
	}

	return baseURL, uri, nil
}

func (cli *Client) getContext(opt RequestOptions) (context.Context, context.CancelFunc) {
	if opt.Context != nil {
		return opt.Context, nil
	}

	baseCtx := context.TODO()

	timeout := RequestTimeoutDefault
	if opt.Timeout > 0 || opt.Timeout == RequestTimeoutNone {
		timeout = opt.Timeout
	} else if cli.timeout > 0 || cli.timeout == RequestTimeoutNone {
		timeout = cli.timeout
	}

	if timeout == RequestTimeoutNone {
		return context.WithCancel(baseCtx)
	}

	return context.WithTimeout(baseCtx, time.Duration(timeout)*time.Millisecond)
}
