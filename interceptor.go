package request

import "net/http"

// RequestInterceptor is a function to intercept the requests.
type RequestInterceptor func(*http.Request) error

// ResponseInterceptor is a function to intercept the responses.
type ResponseInterceptor func(*http.Response) error

// doRequestIntercept executes the request interceptors, it'll terminate if the interceptor
// function returns an error.
func (cli *Client) doRequestIntercept(req *http.Request) error {
	interceptors := cli.reqInterceptors
	if len(interceptors) == 0 {
		return nil
	}

	for _, interceptor := range interceptors {
		err := interceptor(req)
		if err != nil {
			return err
		}
	}

	return nil
}

// doResponseIntercept executes the response interceptors, it'll terminate if the interceptor
// function returns an error.
func (cli *Client) doResponseIntercept(resp *http.Response) error {
	interceptors := cli.respInterceptors
	if len(interceptors) == 0 {
		return nil
	}

	for _, interceptor := range interceptors {
		err := interceptor(resp)
		if err != nil {
			return err
		}
	}

	return nil
}
