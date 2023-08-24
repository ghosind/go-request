package request

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/ghosind/go-assert"
)

func TestSetBasicAuthChain(t *testing.T) {
	a := assert.New(t)

	data, _, err := ToObject[testResponse](Req("http://localhost:8080").
		SetBasicAuth("user", "pass").
		Do())
	a.NilNow(err)
	a.NotNilNow(data)
	a.NotNilNow(data.Token)
	a.EqualNow(*data.Token, "Basic dXNlcjpwYXNz")
}

func TestSetBaseURLChain(t *testing.T) {
	a := assert.New(t)

	data, _, err := ToObject[testResponse](Req("/test").
		SetBaseURL("http://localhost:8080").
		Do())
	a.NilNow(err)
	a.NotNilNow(data)
	a.NotNilNow(data.Path)
	a.EqualNow(*data.Path, "/test")
}

func TestSetBodyChain(t *testing.T) {
	a := assert.New(t)

	data, _, err := ToObject[testResponse](Req("http://localhost:8080").
		POST().
		SetBody(map[string]any{
			"age":      18,
			"greeting": "hello",
		}).
		Do())
	a.NilNow(err)
	a.NotNilNow(data)
	a.NotNilNow(data.Body)
	a.EqualNow(*data.Body, `{"age":18,"greeting":"hello"}`)
}

func TestSetContentTypeChain(t *testing.T) {
	a := assert.New(t)

	data, _, err := ToObject[testResponse](Req("http://localhost:8080").
		POST().
		SetContentType(RequestContentTypeJSON).
		Do())
	a.NilNow(err)
	a.NotNilNow(data)
	a.NotNilNow(data.Headers)
	a.EqualNow((*data.Headers)["Content-Type"], []string{"application/json"})

	_, err = Req("http://localhost:8080").
		POST().
		SetContentType("unknown").
		Do()
	a.NotNilNow(err)
}

func TestSetContextChain(t *testing.T) {
	a := assert.New(t)

	ctx := context.Background()
	_, err := Req("http://localhost:8080").SetContext(ctx).Do()
	a.NilNow(err)

	ctx, canFunc := context.WithCancel(context.Background())
	go func() {
		_, err = Req("http://localhost:8080").SetContext(ctx).Do()
		a.NotNilNow(err)
	}()
	canFunc()
}

func TestSetDisableDecompressChain(t *testing.T) {
	a := assert.New(t)

	resp, err := Req("http://localhost:8080").
		AddHeader("Accept-Encoding", "gzip").
		POST().
		Do()
	a.NilNow(err)
	a.NotNilNow(resp)
	a.EqualNow(resp.Header.Get("Content-Encoding"), "")

	resp, err = Req("http://localhost:8080").
		AddHeader("Accept-Encoding", "gzip").
		POST().
		SetDisableDecompress(true).
		Do()
	a.NilNow(err)
	a.NotNilNow(resp)
	a.EqualNow(resp.Header.Get("Content-Encoding"), "gzip")
}

func TestSetHeadersChain(t *testing.T) {
	a := assert.New(t)

	data, _, err := ToObject[testResponse](Req("http://localhost:8080").
		AddHeader("TEST-1", "test").
		SetHeaders(map[string][]string{
			"TEST-1": {"hello"},
			"TEST-2": {"123", "456"},
		}).
		Do())
	a.NilNow(err)
	a.NotNilNow(data)
	a.NotNilNow(data.Headers)
	a.Equal((*data.Headers)["Test-1"], []string{"hello"})
	a.Equal((*data.Headers)["Test-2"], []string{"123", "456"})

	data, _, err = ToObject[testResponse](Req("http://localhost:8080").
		AddHeader("TEST", "0").
		SetHeader("TEST", []string{"123", "456"}).
		Do())
	a.NilNow(err)
	a.NotNilNow(data)
	a.NotNilNow(data.Headers)
	a.Equal((*data.Headers)["Test"], []string{"123", "456"})

	data, _, err = ToObject[testResponse](Req("http://localhost:8080").
		SetHeader("TEST", []string{"123"}).
		AddHeader("TEST", "456").
		Do())
	a.NilNow(err)
	a.NotNilNow(data)
	a.NotNilNow(data.Headers)
	a.Equal((*data.Headers)["Test"], []string{"123", "456"})
}

func TestSetMaxRedirectsChain(t *testing.T) {
	a := assert.New(t)

	resp, err := Req("http://localhost:8080/redirect").SetMaxRedirects(3).Do()
	a.NilNow(err)

	locationUrl := resp.Header.Get("Location")
	location, err := url.Parse(locationUrl)
	a.NilNow(err)

	tried := location.Query().Get("tried")
	a.EqualNow(tried, "3")
}

func TestSetMethodChain(t *testing.T) {
	a := assert.New(t)

	for _, method := range []string{
		http.MethodDelete,
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut,
	} {
		data, _, err := ToObject[testResponse](Req("http://localhost:8080").SetMethod(method).Do())
		a.NilNow(err, method)

		if method == "HEAD" {
			// no body for HEAD
			continue
		}

		a.NotNilNow(data.Method)
		a.EqualNow(*data.Method, method)
	}
}

func TestMethodsChain(t *testing.T) {
	a := assert.New(t)

	req := Req("http://localhost:8080")

	for method, fn := range map[string]func() *RequestOptions{
		"DELETE":  req.DELETE,
		"GET":     req.GET,
		"HEAD":    req.HEAD,
		"OPTIONS": req.OPTIONS,
		"PATCH":   req.PATCH,
		"POST":    req.POST,
		"PUT":     req.PUT,
	} {
		data, _, err := ToObject[testResponse](fn().Do())
		a.NilNow(err, method)

		if method == "HEAD" {
			// no body for HEAD
			continue
		}

		a.NotNilNow(data.Method)
		a.EqualNow(*data.Method, method)
	}
}

func TestSetParametersChain(t *testing.T) {
	a := assert.New(t)

	data, _, err := ToObject[testResponse](Req("http://localhost:8080").
		AddParameter("status", "-10").
		SetParameters(map[string][]string{
			"text":   {"test"},
			"status": {"0", "10"},
		}).
		Do())
	a.NilNow(err)
	a.NotNilNow(data)
	a.NotNilNow(data.UserAgent)
	a.Equal(*data.Query, "status=0&status=10&text=test")

	data, _, err = ToObject[testResponse](Req("http://localhost:8080").
		AddParameter("status", "-10").
		SetParameter("status", []string{"0", "10"}).
		Do())
	a.NilNow(err)
	a.NotNilNow(data)
	a.NotNilNow(data.UserAgent)
	a.Equal(*data.Query, "status=0&status=10")

	data, _, err = ToObject[testResponse](Req("http://localhost:8080").
		SetParameter("status", []string{"0"}).
		AddParameter("status", "10").
		Do())
	a.NilNow(err)
	a.NotNilNow(data)
	a.NotNilNow(data.UserAgent)
	a.Equal(*data.Query, "status=0&status=10")
}

func TestSetTimeoutChain(t *testing.T) {
	a := assert.New(t)

	_, err := Req("http://localhost:8080").SetTimeout(RequestTimeoutNoLimit).Do()
	a.NilNow(err)

	// TODO: 1ms timeout is not stable to test
	// _, err = Req("http://localhost:8080").SetTimeout(1).Do()
	// a.NotNilNow(err)
}

func TestSetUserAgentChain(t *testing.T) {
	a := assert.New(t)

	data, _, err := ToObject[testResponse](Req("http://localhost:8080").
		Do())
	a.NilNow(err)
	a.NotNilNow(data)
	a.NotNilNow(data.UserAgent)
	a.Equal(*data.UserAgent, RequestDefaultUserAgent)

	data, _, err = ToObject[testResponse](Req("http://localhost:8080").
		SetUserAgent("test-bot").
		Do())
	a.NilNow(err)
	a.NotNilNow(data)
	a.NotNilNow(data.UserAgent)
	a.Equal(*data.UserAgent, "test-bot")
}

func TestSetValidateStatusChain(t *testing.T) {
	a := assert.New(t)

	validateStatus := func(code int) bool {
		return code == 400
	}

	resp, err := Req("http://localhost:8080/").
		SetValidateStatus(validateStatus).
		Do()
	a.NotNilNow(err)
	a.NotNilNow(resp)
	a.EqualNow(resp.StatusCode, 200)

	resp, err = Req("http://localhost:8080/status").
		AddParameter("status", "400").
		SetValidateStatus(validateStatus).
		Do()
	a.NilNow(err)
	a.NotNilNow(resp)
	a.EqualNow(resp.StatusCode, 400)
}

func TestSetURLChain(t *testing.T) {
	a := assert.New(t)

	data, _, err := ToObject[testResponse](Req("/").
		SetBaseURL("http://example.com").
		SetURL("http://localhost:8080/test").
		Do(),
	)
	a.NilNow(err)
	a.NotNilNow(data.Path)
	a.EqualNow(*data.Path, "/test")
}

func TestRequestOptionsDo(t *testing.T) {
	a := assert.New(t)

	resp, err := Req("http://localhost:8080").Do()
	a.NilNow(err)
	a.EqualNow(resp.StatusCode, 200)

	resp, err = New().Req("http://localhost:8080").Do()
	a.NilNow(err)
	a.EqualNow(resp.StatusCode, 200)

	resp, err = (&RequestOptions{}).SetURL("http://localhost:8080").Do()
	a.NilNow(err)
	a.EqualNow(resp.StatusCode, 200)
}
