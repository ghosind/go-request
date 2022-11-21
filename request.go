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
}

func Request(url string, opts ...RequestOptions) (*http.Response, error) {
	return defaultClient.Request(url, opts...)
}

func (cli *Client) Request(url string, opts ...RequestOptions) (*http.Response, error) {
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
	defer canFunc()

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

	// TODO: headers

	return req, canFunc, nil
}

func (cli *Client) parseURL(u string, opt RequestOptions) (string, error) {
	base := ""
	if opt.BaseURL != "" {
		base = opt.BaseURL
	} else if cli.BaseURL != "" {
		base = cli.BaseURL
	}

	if base == "" {
		base = u
		u = ""
	}

	obj, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	if u != "" {
		obj.Path = path.Join(obj.Path, u)
	}

	query := obj.Query()

	for k, vv := range opt.Parameters {
		for _, v := range vv {
			query.Add(k, v)
		}
	}

	obj.RawQuery = query.Encode()

	return obj.String(), nil
}

func (cli *Client) getContext(opt RequestOptions) (context.Context, context.CancelFunc) {
	var baseCtx context.Context

	if opt.Context != nil {
		baseCtx = opt.Context
	} else {
		baseCtx = context.TODO()
	}

	timeout := 0
	if opt.Timeout != 0 {
		timeout = opt.Timeout
	} else if cli.timeout != 0 {
		timeout = int(cli.timeout)
	} else {
		timeout = DefaultTimeout
	}

	if timeout == -1 {
		return context.WithCancel(baseCtx)
	}

	return context.WithTimeout(baseCtx, time.Duration(timeout)*time.Millisecond)
}
