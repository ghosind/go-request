package request

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"
)

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
	defer canFunc()

	resp, err := cli.sendRequestWithInterceptors(req, opt)
	if err != nil {
		return resp, err
	}

	return cli.handleResponse(resp, opt)
}

// sendRequestWithInterceptors tries to execute the request and response interceptors and
// sends the request.
func (cli *Client) sendRequestWithInterceptors(
	req *http.Request,
	opt RequestOptions,
) (*http.Response, error) {
	err := cli.doRequestIntercept(req)
	if err != nil {
		return nil, err
	}

	resp, err := cli.sendRequest(req, opt)
	if err != nil {
		return nil, err
	}

	err = cli.doResponseIntercept(resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// sendRequest gets an HTTP client from the HTTP clients pool and sends the request. It tries to
// re-send the request when it fails to make the request and the number of attempts is less than
// the maximum limitation.
func (cli *Client) sendRequest(req *http.Request, opt RequestOptions) (*http.Response, error) {
	attempt := 0
	maxAttempt := 1
	if opt.MaxAttempt > 0 {
		maxAttempt = opt.MaxAttempt
	}

	httpClient := cli.getHTTPClient(opt)
	defer func() {
		cli.clientPool.Put(httpClient)
	}()

	for {
		attempt++

		resp, err := httpClient.Do(req)
		if err == nil || (err != nil && attempt >= maxAttempt) {
			return resp, err
		} else if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return resp, err
		}
	}
}

// handleResponse handle the response that decompresses the body of the response if it was
// compressed, and validates the status code.
func (cli *Client) handleResponse(
	resp *http.Response,
	opt RequestOptions,
) (*http.Response, error) {
	if !opt.DisableDecompress {
		resp = cli.decodeResponseBody(resp)
	}

	return cli.validateResponse(resp, opt)
}

// decodeResponseBody tries to get the encoding type of the response's content, and decode
// (decompress) it if the response's body was compressed by `gzip` or `deflate`.
func (cli *Client) decodeResponseBody(resp *http.Response) *http.Response {
	switch resp.Header.Get("Content-Encoding") {
	case "deflate":
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return resp
		}

		reader := flate.NewReader(bytes.NewReader(data))
		resp.Body.Close()
		resp.Body = reader
		resp.Header.Del("Content-Encoding")
	case "gzip", "x-gzip":
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return resp
		}
		resp.Body.Close()
		resp.Body = reader
		resp.Header.Del("Content-Encoding")
	}

	return resp
}

// validateResponse validates the status code of the response, and returns fail if the result of
// the validation is false.
func (cli *Client) validateResponse(
	resp *http.Response,
	opt RequestOptions,
) (*http.Response, error) {
	var validateStatus func(int) bool
	if opt.ValidateStatus != nil {
		validateStatus = opt.ValidateStatus
	} else if cli.ValidateStatus != nil {
		validateStatus = cli.ValidateStatus
	} else {
		validateStatus = cli.defaultValidateStatus
	}

	status := resp.StatusCode
	ok := validateStatus(status)
	if !ok {
		return resp, fmt.Errorf("request failed with status code %d", status)
	}

	return resp, nil
}

// makeRequest creates a new `http.Request` object with the specific HTTP method, request url, and
// other configurations.
func (cli *Client) makeRequest(
	method, url string,
	opt RequestOptions,
) (*http.Request, context.CancelFunc, error) {
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

// parseURL gets the URL of the request and adds the parameters into the query of the request.
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

	obj.RawQuery = cli.getQueryParameters(obj.Query(), opt)

	return obj.String(), nil
}

// getQueryParameters get the parameters of the request from the request options and the client's
// parameters.
func (cli *Client) getQueryParameters(query url.Values, opt RequestOptions) string {
	if opt.Parameters != nil {
		for k, vv := range opt.Parameters {
			if !query.Has(k) {
				query[k] = make([]string, 0, len(vv))
			}
			query[k] = append(query[k], vv...)
		}
	}

	if cli.Parameters != nil {
		for k, vv := range cli.Parameters {
			if query.Has(k) {
				continue
			}

			query[k] = make([]string, 0, len(vv))
			query[k] = append(query[k], vv...)
		}
	}

	if opt.ParametersSerializer != nil {
		return opt.ParametersSerializer(query)
	} else if cli.ParametersSerializer != nil {
		return cli.ParametersSerializer(query)
	}

	return query.Encode()
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
		return opt.Context, func() {} // empty cancel function, just do nothing
	}

	baseCtx := context.Background()

	timeout := RequestTimeoutDefault
	if opt.Timeout > 0 || opt.Timeout == RequestTimeoutNoLimit {
		timeout = opt.Timeout
	} else if cli.Timeout > 0 || cli.Timeout == RequestTimeoutNoLimit {
		timeout = cli.Timeout
	}

	if timeout == RequestTimeoutNoLimit {
		return baseCtx, func() {} // empty cancel function, just do nothing
	} else {
		return context.WithTimeout(baseCtx, time.Duration(timeout)*time.Millisecond)
	}
}
