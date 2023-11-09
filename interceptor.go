package request

import "net/http"

// RequestInterceptor is a function to intercept the requests.
type RequestInterceptor func(*http.Request) error

// ResponseInterceptor is a function to intercept the responses.
type ResponseInterceptor func(*http.Response) error

// requestInterceptor is a wrapper object for the request interceptor function and its ID.
type requestInterceptor struct {
	// ID is the ID of the interceptor.
	ID uint64
	// Interceptor is the request intercept function.
	Interceptor RequestInterceptor
}

// responseInterceptor is a wrapper object for the response interceptor function and its ID.
type responseInterceptor struct {
	// ID is the ID of the interceptor.
	ID uint64
	// Interceptor is the response intercept function.
	Interceptor ResponseInterceptor
}

// doRequestIntercept executes the request interceptors, it'll terminate if the interceptor
// function returns an error.
func (cli *Client) doRequestIntercept(req *http.Request) error {
	interceptors := cli.reqInterceptors
	if len(interceptors) == 0 {
		return nil
	}

	cli.interceptorMutex.RLock()
	defer cli.interceptorMutex.RUnlock()

	for _, interceptor := range interceptors {
		err := interceptor.Interceptor(req)
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

	cli.interceptorMutex.RLock()
	defer cli.interceptorMutex.RUnlock()

	for _, interceptor := range interceptors {
		err := interceptor.Interceptor(resp)
		if err != nil {
			return err
		}
	}

	return nil
}
