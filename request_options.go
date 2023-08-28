package request

import (
	"context"
	"net/http"
)

// RequestOptions is the config for a request.
type RequestOptions struct {
	// Auth indicates that HTTP Basic auth should be used. It will set an `Authorization` header,
	// and it'll also overwriting any existing `Authorization` field in the request header.
	//
	//	resp, err := request.Request("https://example.com", request.RequestOptions{
	//	  Auth: &request.BasicAuthConfig{
	//	    Username: "user",
	//	    Password: "pass",
	//	  },
	//	})
	Auth *BasicAuthConfig
	// BaseURL will prepended to the url of the request unless the url is absolute.
	//
	//	resp, err := request.Request("/test", request.RequestOptions{
	//	  BaseURL: "http://example.com",
	//	})
	//	// http://example.com/test
	BaseURL string
	// Body is the data to be sent as the request body. It'll be encoded with the content type
	// specified by the `ContentType` field in the request options, or encoded as a JSON if the
	// `ContentType` field is empty. It'll skip the encode processing if the value is a string or a
	// slice of bytes.
	//
	//	resp, err := request.Request("http://example.com", request.RequestOptions{
	//	  Method: http.MethodPost,
	//	  Body: "Hello world!", // with raw string
	//	})
	//
	//	resp, err := request.Request("http://example.com", request.RequestOptions{
	//	  Method: http.MethodPost,
	//	  // with struct/map, and it'll encoding by the value of ContentType field.
	//	  Body: map[string]any{
	//	    "data": "Hello world",
	//	  },
	//	})
	Body any
	// ContentType indicates the type of data that will encode and send to the server. Available
	// options are: "json", default "json".
	//
	//	request.POST("http://example.com", request.RequestOptions{
	//	  ContentType: request.RequestContentTypeJSON, // "json"
	//	  // ...
	//	})
	ContentType string
	// Context is a `context.Content` object that is used for manipulating the request by yourself.
	// The `Timeout` field will be ignored if this value is not empty, and you need to control
	// timeout by yourself.
	//
	//	ctx, canFunc := context.WithTimeout(context.Background(), time.Second)
	//	defer canFunc()
	//	resp, err := request.Request("http://example.com", request.RequestOptions{
	//	  Context: ctx,
	//	})
	Context context.Context
	// DisableDecompress indicates whether or not disable decompression of the response body
	// automatically. If it is set to `true`, it will not decompress the response body.
	DisableDecompress bool
	// Headers are custom headers to be sent.
	//
	//	resp, err := request.Request("http://example.com", request.RequestOptions{
	//	  Headers: map[string][]string{
	//	    "Authorization": {"Bearer XXXXX"},
	//	  },
	//	})
	Headers map[string][]string
	// MaxRedirects defines the maximum number of redirects, default 5.
	MaxRedirects int
	// Method indicates the HTTP method of the request, default GET.
	//
	//	request.Request("http://example.com", request.RequestOptions{
	//	  Method: http.MethodPost, // "POST"
	//	})
	Method string
	// Parameters are the URL parameters to be sent with the request.
	//
	//	resp, err := request.Request("http://example.com", request.RequestOptions{
	//	  Parameters: map[string][]string{
	//	    "name": {"John"},
	//	  },
	//	})
	//	// http://example.com?name=John
	Parameters map[string][]string
	// ParametersSerializer is a function to charge of serializing the URL query parameters.
	ParametersSerializer func(map[string][]string) string
	// Timeout specifies the number of milliseconds before the request times out. This value will be
	// ignored if the `Content` field in the request options is set. It indicates no time-out
	// limitation if the value is -1.
	Timeout int
	// UserAgent sets the client's User-Agent field in the request header. It'll overwrite the value
	// of the `User-Agent` field in the request headers.
	UserAgent string
	// ValidateStatus defines whether the status code of the response is valid or not, and it'll
	// return an error if fails to validate the status code. Default, it sets the result to fail if
	// the status code is less than 200, or greater than and equal to 400.
	//
	//	resp, err := request.Request("http://example.com", request.RequestOptions{
	//	  ValidateStatus: func (status int) bool {
	//	    // Only success if the status code of response is 2XX
	//	    return status >= http.StatusOk && status <= http.StatusMultipleChoices
	//	  },
	//	})
	ValidateStatus func(int) bool

	// client is the Client instance of the request. This field is for chaining API only, and it'll
	// initialized by the `Req` method. If you create a `requestOptions` manually, the request will
	// use the default client.
	client *Client
	// url is the destination URL of the request, and this field if for chaining API only.
	url string
}

// BasicAuthConfig indicates the config of the HTTP Basic Auth that is used for the request.
type BasicAuthConfig struct {
	// Username indicates the username used for HTTP Basic Auth
	Username string
	// Password indicates the password used for HTTP Basic Auth
	Password string
}

// SetBasicAuth sets the username and the password as the HTTP Basic Auth to the request.
//
//	request.Req("http://example.com").
//	  SetBasicAuth("user", "pass").
//	  Do()
func (opt *RequestOptions) SetBasicAuth(username, password string) *RequestOptions {
	if opt.Auth == nil {
		opt.Auth = new(BasicAuthConfig)
	}

	opt.Auth.Username = username
	opt.Auth.Password = password

	return opt
}

// SetBaseURL sets the base URL of the request.
//
//	request.Req("/test").
//	  SetBaseURL("http://example.com").
//	  Do()
//	// http://example.com/test
func (opt *RequestOptions) SetBaseURL(baseURL string) *RequestOptions {
	opt.BaseURL = baseURL

	return opt
}

// SetBody sets the request body to the request.
//
//	request.Req("http://example.com").
//	  POST().
//	  SetBody(map[string]string{ "greeting": "Hello world!" }).
//	  Do()
func (opt *RequestOptions) SetBody(body any) *RequestOptions {
	opt.Body = body

	return opt
}

// SetContentType sets the encoding type of the request content, default `json`.
//
//	request.Req("http://example.com").
//	  POST().
//	  SetContentType(request.RequestContentTypeJSON).
//	  SetBody(map[string]string{ "greeting": "Hello world!" }).
//	  Do()
func (opt *RequestOptions) SetContentType(contentType string) *RequestOptions {
	opt.ContentType = contentType

	return opt
}

// SetContext sets the context of the request, and it will skip the value of the `Timeout` field in
// the request options if context is not empty.
//
//	ctx, canFunc := context.WithTimeout(context.Background, time.Second)
//	defer canFunc()
//
//	request.Req("http://example.com").
//	  SetContext(ctx).
//	  Do()
func (opt *RequestOptions) SetContext(ctx context.Context) *RequestOptions {
	opt.Context = ctx

	return opt
}

// SetDisableDecompress sets whether to decompress the response body or not, and sets true to
// disable automatic decompression.
//
//	request.Req("http://example.com").
//	  AddHeader("Accept-Encoding", "gzip").
//	  SetDisableDecompress(true).
//	  Do()
func (opt *RequestOptions) SetDisableDecompress(isDisable bool) *RequestOptions {
	opt.DisableDecompress = isDisable

	return opt
}

// AddHeader add the value to the request headers with the specified key.
//
//	request.Req("http://example.com").
//	  AddHeader("Content-Type", []string{"application/json"}).
//	  Do()
func (opt *RequestOptions) AddHeader(key, val string) *RequestOptions {
	if opt.Headers == nil {
		opt.Headers = make(map[string][]string)
	}

	values, ok := opt.Headers[key]
	if ok {
		opt.Headers[key] = append(values, val)
	} else {
		opt.Headers[key] = []string{val}
	}

	return opt
}

// SetHeader sets the value of the request headers with the specified key, and it will overwrite
// the value of the key.
//
//	request.Req("http://example.com").
//	  SetHeader("Content-Type", []string{"application/json"}).
//	  Do()
func (opt *RequestOptions) SetHeader(key string, values []string) *RequestOptions {
	if opt.Headers == nil {
		opt.Headers = make(map[string][]string)
	}

	opt.Headers[key] = values

	return opt
}

// SetHeaders sets the headers of the request, and it'll overwrite all the headers if it was set
// before.
//
//	request.Req("http://example.com").
//	  SetHeaders(map[string][]string{
//	    "Authorization": {token},
//	    "Accept-Encoding": {"gzip"},
//	  }).
//	  Do()
func (opt *RequestOptions) SetHeaders(headers map[string][]string) *RequestOptions {
	opt.Headers = headers

	return opt
}

// SetMaxRedirects sets the maximum number of redirects for the request.
//
//	Req("http://example.com").
//	  SetMaxRedirects(3).
//	  Do()
func (opt *RequestOptions) SetMaxRedirects(maxRedirects int) *RequestOptions {
	opt.MaxRedirects = maxRedirects

	return opt
}

// SetMethod sets the HTTP method of the request.
//
//	Req("http://localhost:8080").
//	  SetMethod(http.MethodPost).
//	  Do()
//	// POST("http://localhost:8080")
func (opt *RequestOptions) SetMethod(method string) *RequestOptions {
	opt.Method = method

	return opt
}

// DELETE sets the HTTP method of the request to `DELETE`.
//
//	Req("http://example.com").
//	  DELETE().
//	  Do()
//	// DELETE http://example.com
func (opt *RequestOptions) DELETE() *RequestOptions {
	opt.Method = http.MethodDelete

	return opt
}

// GET sets the HTTP method of the request to `GET`.
//
//	Req("http://example.com").
//	  GET().
//	  Do()
//	// GET http://example.com
func (opt *RequestOptions) GET() *RequestOptions {
	opt.Method = http.MethodGet

	return opt
}

// HEAD sets the HTTP method of the request to `HEAD`.
//
//	Req("http://example.com").
//	  HEAD().
//	  Do()
//	// HEAD http://example.com
func (opt *RequestOptions) HEAD() *RequestOptions {
	opt.Method = http.MethodHead

	return opt
}

// OPTIONS sets the HTTP method of the request to `OPTIONS`.
//
//	Req("http://example.com").
//	  OPTIONS().
//	  Do()
//	// OPTIONS http://example.com
func (opt *RequestOptions) OPTIONS() *RequestOptions {
	opt.Method = http.MethodOptions

	return opt
}

// PATCH sets the HTTP method of the request to `PATCH`.
//
//	Req("http://example.com").
//	  PATCH().
//	  Do()
//	// PATCH http://example.com
func (opt *RequestOptions) PATCH() *RequestOptions {
	opt.Method = http.MethodPatch

	return opt
}

// POST sets the HTTP method of the request to `POST`.
//
//	Req("http://example.com").
//	  POST().
//	  Do()
//	// POST http://example.com
func (opt *RequestOptions) POST() *RequestOptions {
	opt.Method = http.MethodPost

	return opt
}

// PUT sets the HTTP method of the request to `PUT`.
//
//	Req("http://example.com").
//	  PUT().
//	  Do()
//	// PUT http://example.com
func (opt *RequestOptions) PUT() *RequestOptions {
	opt.Method = http.MethodPut

	return opt
}

// SetParameter add the value to the query parameter with the specified key.
//
//	request.Req("http://example.com").
//	  AddParameter("status", "0").
//	  AddParameter("status", "10").
//	  Do()
//	// http://example.com?status=0&status=10
func (opt *RequestOptions) AddParameter(key, val string) *RequestOptions {
	if opt.Parameters == nil {
		opt.Parameters = make(map[string][]string)
	}

	values, ok := opt.Parameters[key]
	if ok {
		opt.Parameters[key] = append(values, val)
	} else {
		opt.Parameters[key] = []string{val}
	}

	return opt
}

// SetParameter sets the value of the query parameter with the specified key, and it will overwrite
// the value of the key.
//
//	request.Req("http://example.com").
//	  SetParameter("status", []string{"0", "-10"}).
//	  SetParameter("status", []string{"0", "10"}). // it'll overwrite previous one
//	  Do()
//	// http://example.com?status=0&status=10
func (opt *RequestOptions) SetParameter(key string, values []string) *RequestOptions {
	if opt.Parameters == nil {
		opt.Parameters = make(map[string][]string)
	}

	opt.Parameters[key] = values

	return opt
}

// SetParameters sets the query parameters of the request, and it'll overwrite all the parameters
// if it was set before.
//
//	request.Req("http://example.com").
//	  SetParameters(map[string][]string{
//	    "text": {"test"},
//	    "status": {"0","10"},
//	  }).
//	  Do()
//	// http://example.com?text=test&status=0&status=10
func (opt *RequestOptions) SetParameters(parameters map[string][]string) *RequestOptions {
	opt.Parameters = parameters

	return opt
}

// SetParametersSerializer set the function that to charge of serializing the URL query parameters.
func (opt *RequestOptions) SetParametersSerializer(
	fn func(map[string][]string) string,
) *RequestOptions {
	opt.ParametersSerializer = fn

	return opt
}

// SetTimeout sets the timeout of the request in the milliseconds.
//
//	request.Req("http://example.com").
//	  SetTimeout(3000). // Timeout is 3 seconds
//	  Do()
func (opt *RequestOptions) SetTimeout(timeout int) *RequestOptions {
	opt.Timeout = timeout

	return opt
}

// SetUserAgent sets the value of the `User-Agent` field in the request headers.
//
//	request.Req("http://example.com").
//	  SetUserAgent("test-bot/1.0").
//	  Do()
func (opt *RequestOptions) SetUserAgent(userAgent string) *RequestOptions {
	opt.UserAgent = userAgent

	return opt
}

// SetValidateStatus sets the function that defines whether the status code of the response is
// valid or not.
//
//	request.Req("http://example.com").
//	  // only status codes 2XX are allowed
//	  SetValidateStatus(func (code int) bool { return code >= 200 && code < 300 }).
//	  Do()
func (opt *RequestOptions) SetValidateStatus(validateStatus func(int) bool) *RequestOptions {
	opt.ValidateStatus = validateStatus

	return opt
}

// SetURL overwrites the URL of the request.
//
//	request.Req("http://localhost").
//	  SetURL("http://example.com").
//	  Do()
//	// GET http://example.com
func (opt *RequestOptions) SetURL(url string) *RequestOptions {
	opt.url = url

	return opt
}

// Do makes the request by the request options, and returns the response or the error.
//
//	request.Req("http://example.com").
//	  Do()
//	// GET("http://example.com")
func (opt *RequestOptions) Do() (*http.Response, error) {
	cli := opt.client
	if cli == nil {
		cli = defaultClient
	}

	return cli.request(opt.Method, opt.url, *opt)
}
