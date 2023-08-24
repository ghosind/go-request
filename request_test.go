package request

import (
	"context"
	"math"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ghosind/go-assert"
)

func TestMakeRequest(t *testing.T) {
	a := assert.New(t)
	cli := New()

	// no error
	_, _, err := cli.makeRequest("", "http://example.com", RequestOptions{})
	a.NilNow(err)

	// invalid HTTP method
	_, _, err = cli.makeRequest("UNKNOWN", "http://example.com", RequestOptions{})
	a.NotNilNow(err)

	// invalid url
	_, _, err = cli.makeRequest("", "", RequestOptions{})
	a.NotNilNow(err)

	// invalid content type for encoding body
	_, _, err = cli.makeRequest("", "http://example.com", RequestOptions{
		Body:        []string{"Test"},
		ContentType: "unknown",
	})
	a.NotNilNow(err)

	// invalid content type for headers
	_, _, err = cli.makeRequest("", "http://example.com", RequestOptions{
		ContentType: "unknown",
	})
	a.NotNilNow(err)
}

func TestGetRequestMethod(t *testing.T) {
	a := assert.New(t)
	cli := New()

	method, err := cli.getRequestMethod("")
	a.Nil(err)
	a.Equal(method, "GET")

	// valid methods
	for _, method := range []string{"Connect", "delete", "get", http.MethodHead, "Options", "PATCH", "PoST", "PuT", "TRACE"} {
		ret, err := cli.getRequestMethod(method)
		a.Nil(err)
		a.Equal(ret, strings.ToUpper(method))
	}

	_, err = cli.getRequestMethod("UNKNOWN")
	a.NotNil(err)
}

func TestAttachRequestHeaders(t *testing.T) {
	a := assert.New(t)
	cli := New()

	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		a.Fatalf("Unexpected error: %v", err)
	}

	// No error
	err = cli.attachRequestHeaders(req, RequestOptions{})
	a.NilNow(err)
	req.Header.Del("Content-Type")

	// invalid content type
	err = cli.attachRequestHeaders(req, RequestOptions{ContentType: "unknown"})
	a.NotNilNow(err)
}

func TestSetHeaders(t *testing.T) {
	a := assert.New(t)
	cli := New(Config{
		Headers: map[string][]string{
			"Key1": {"V1"},
			"Key2": {"V2"},
		},
	})

	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		a.Fatalf("Unexpected error: %v", err)
	}

	cli.setHeaders(req, RequestOptions{
		Headers: map[string][]string{
			"Key2": {"V1"},
			"Key3": {"V3"},
		},
	})
	a.Equal(req.Header.Get("Key1"), "V1")
	a.Equal(req.Header.Get("Key2"), "V1")
	a.Equal(req.Header.Get("Key3"), "V3")
}

func TestBasicAuth(t *testing.T) {
	a := assert.New(t)
	cli := New()

	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		a.Fatalf("Unexpected error: %v", err)
	}

	err = cli.attachRequestHeaders(req, RequestOptions{})
	a.NilNow(err)
	a.Equal(req.Header.Get("Authorization"), "")
	req.Header.Del("Authorization")

	err = cli.attachRequestHeaders(req, RequestOptions{
		Auth: &BasicAuthConfig{
			Username: "user",
			Password: "pass",
		},
	})
	a.NilNow(err)
	a.Equal(req.Header.Get("Authorization"), "Basic dXNlcjpwYXNz")
}

func TestSetContentType(t *testing.T) {
	a := assert.New(t)
	cli := New()

	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		a.Fatalf("Unexpected error: %v", err)
	}

	err = cli.setContentType(req, RequestOptions{})
	a.NilNow(err)
	a.Equal(req.Header.Get("Content-Type"), "application/json")
	req.Header.Del("Content-Type")

	err = cli.setContentType(req, RequestOptions{
		ContentType: "json",
	})
	a.NilNow(err)
	a.Equal(req.Header.Get("Content-Type"), "application/json")
	req.Header.Del("Content-Type")

	req.Header.Set("Content-Type", "application/vnd.github+json")
	err = cli.setContentType(req, RequestOptions{
		ContentType: "json",
	})
	a.NilNow(err)
	a.Equal(req.Header.Get("Content-Type"), "application/vnd.github+json")
	req.Header.Del("Content-Type")

	err = cli.setContentType(req, RequestOptions{
		ContentType: "unknown",
	})
	a.NotNilNow(err)
}

func TestSetUserAgent(t *testing.T) {
	a := assert.New(t)
	cli := New()

	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		a.Fatalf("Unexpected error: %v", err)
	}

	cli.setUserAgent(req, RequestOptions{})
	a.Equal(req.Header.Get("User-Agent"), RequestDefaultUserAgent)
	req.Header.Del("User-Agent")

	cli.UserAgent = "Test-HTTP-Client"
	cli.setUserAgent(req, RequestOptions{})
	a.Equal(req.Header.Get("User-Agent"), "Test-HTTP-Client")
	req.Header.Del("User-Agent")

	cli.setUserAgent(req, RequestOptions{UserAgent: "Test-Client"})
	a.Equal(req.Header.Get("User-Agent"), "Test-Client")
	req.Header.Del("User-Agent")
}

func TestParseURL(t *testing.T) {
	a := assert.New(t)
	cli := New()

	_, err := cli.parseURL("", RequestOptions{})
	a.NotNilNow(err)

	_, err = cli.parseURL("some invalid url", RequestOptions{})
	a.NotNilNow(err)

	url, err := cli.parseURL("http://example.com", RequestOptions{})
	a.NilNow(err)
	a.Equal(url, "http://example.com")

	url, err = cli.parseURL("test", RequestOptions{
		BaseURL: "http://example.com",
	})
	a.NilNow(err)
	a.Equal(url, "http://example.com/test")

	url, err = cli.parseURL("http://example.com?q=test1&w=1", RequestOptions{
		Parameters: map[string][]string{
			"q": {"test2"},
			"t": {"2"},
		},
	})
	a.NilNow(err)
	a.Equal(url, "http://example.com?q=test1&q=test2&t=2&w=1")

	cli = New(Config{
		Parameters: map[string][]string{
			"q": {"test1"},
			"w": {"1"},
		},
	})

	url, err = cli.parseURL("http://example.com", RequestOptions{
		Parameters: map[string][]string{
			"q": {"test2"},
			"t": {"2"},
		},
	})
	a.NilNow(err)
	a.Equal(url, "http://example.com?q=test2&t=2&w=1")
}

func TestGetURL(t *testing.T) {
	a := assert.New(t)
	cli := New()

	_, _, err := cli.getURL("", RequestOptions{})
	a.NotNilNow(err)

	baseUrl, url, err := cli.getURL("http://www.example.com", RequestOptions{})
	a.NilNow(err)
	a.EqualNow(baseUrl, "http://www.example.com")
	a.EqualNow(url, "")

	baseUrl, url, err = cli.getURL("www.example.com", RequestOptions{})
	a.NilNow(err)
	a.EqualNow(baseUrl, "https://www.example.com")
	a.EqualNow(url, "")

	baseUrl, url, err = cli.getURL("", RequestOptions{
		BaseURL: "http://www.example.com",
	})
	a.NilNow(err)
	a.EqualNow(baseUrl, "http://www.example.com")
	a.EqualNow(url, "")

	baseUrl, url, err = cli.getURL("/test", RequestOptions{
		BaseURL: "http://www.example.com",
	})
	a.NilNow(err)
	a.EqualNow(baseUrl, "http://www.example.com")
	a.EqualNow(url, "/test")

	cli.BaseURL = "http://www.example.com"
	baseUrl, url, err = cli.getURL("", RequestOptions{})
	a.NilNow(err)
	a.EqualNow(baseUrl, "http://www.example.com")
	a.EqualNow(url, "")

	baseUrl, url, err = cli.getURL("", RequestOptions{
		BaseURL: "http://www.another.com",
	})
	a.NilNow(err)
	a.EqualNow(baseUrl, "http://www.another.com")
	a.EqualNow(url, "")
}

func TestGetContext(t *testing.T) {
	a := assert.New(t)
	cli := New()

	baseCtx := context.Background()
	ctx, _ := cli.getContext(RequestOptions{
		Context: baseCtx,
	})
	a.Equal(ctx, baseCtx)

	ctx, _ = cli.getContext(RequestOptions{})
	deadline, ok := ctx.Deadline()
	a.EqualNow(ok, true)
	a.TrueNow(math.Abs(float64(1000-(deadline.UnixMilli()-time.Now().UnixMilli()))) < 10)

	ctx, _ = cli.getContext(RequestOptions{
		Timeout: 3000,
	})
	deadline, ok = ctx.Deadline()
	a.EqualNow(ok, true)
	a.TrueNow(math.Abs(float64(3000-(deadline.UnixMilli()-time.Now().UnixMilli()))) < 10)

	ctx, _ = cli.getContext(RequestOptions{
		Timeout: RequestTimeoutNoLimit,
	})
	_, ok = ctx.Deadline()
	a.EqualNow(ok, false)

	cli.Timeout = 3000
	ctx, _ = cli.getContext(RequestOptions{})
	deadline, ok = ctx.Deadline()
	a.EqualNow(ok, true)
	a.TrueNow(math.Abs(float64(3000-(deadline.UnixMilli()-time.Now().UnixMilli()))) < 10)
}
