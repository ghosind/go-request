package request

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ghosind/go-assert"
)

func TestSetContentType(t *testing.T) {
	a := assert.New(t)
	cli := New()

	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		a.Fatalf("Unexpected error: %v", err)
	}

	err = cli.setContentType(req, RequestOptions{})
	a.NilNow(err)
	a.DeepEqual(req.Header.Get("Content-Type"), "application/json")
	req.Header.Del("Content-Type")

	err = cli.setContentType(req, RequestOptions{
		ContentType: "json",
	})
	a.NilNow(err)
	a.DeepEqual(req.Header.Get("Content-Type"), "application/json")
	req.Header.Del("Content-Type")

	req.Header.Set("Content-Type", "application/vnd.github+json")
	err = cli.setContentType(req, RequestOptions{
		ContentType: "json",
	})
	a.NilNow(err)
	a.DeepEqual(req.Header.Get("Content-Type"), "application/vnd.github+json")
	req.Header.Del("Content-Type")

	err = cli.setContentType(req, RequestOptions{
		ContentType: "unknown",
	})
	a.NotNilNow(err)
}

func TestGetRequestMethod(t *testing.T) {
	a := assert.New(t)
	cli := New()

	method, err := cli.getRequestMethod("")
	a.Nil(err)
	a.DeepEqual(method, "GET")

	// valid methods
	for _, method := range []string{"Connect", "delete", "get", http.MethodHead, "Options", "PATCH", "PoST", "PuT", "TRACE"} {
		ret, err := cli.getRequestMethod(method)
		a.Nil(err)
		a.DeepEqual(ret, strings.ToUpper(method))
	}

	_, err = cli.getRequestMethod("UNKNOWN")
	a.NotNil(err)
}

func TestGetURL(t *testing.T) {
	a := assert.New(t)
	cli := New()

	_, _, err := cli.getURL("", RequestOptions{})
	a.NotNilNow(err)

	baseUrl, url, err := cli.getURL("http://www.example.com", RequestOptions{})
	a.NilNow(err)
	a.DeepEqualNow(baseUrl, "http://www.example.com")
	a.DeepEqualNow(url, "")

	baseUrl, url, err = cli.getURL("www.example.com", RequestOptions{})
	a.NilNow(err)
	a.DeepEqualNow(baseUrl, "https://www.example.com")
	a.DeepEqualNow(url, "")

	baseUrl, url, err = cli.getURL("", RequestOptions{
		BaseURL: "http://www.example.com",
	})
	a.NilNow(err)
	a.DeepEqualNow(baseUrl, "http://www.example.com")
	a.DeepEqualNow(url, "")

	baseUrl, url, err = cli.getURL("/test", RequestOptions{
		BaseURL: "http://www.example.com",
	})
	a.NilNow(err)
	a.DeepEqualNow(baseUrl, "http://www.example.com")
	a.DeepEqualNow(url, "/test")

	cli.BaseURL = "http://www.example.com"
	baseUrl, url, err = cli.getURL("", RequestOptions{})
	a.NilNow(err)
	a.DeepEqualNow(baseUrl, "http://www.example.com")
	a.DeepEqualNow(url, "")

	baseUrl, url, err = cli.getURL("", RequestOptions{
		BaseURL: "http://www.another.com",
	})
	a.NilNow(err)
	a.DeepEqualNow(baseUrl, "http://www.another.com")
	a.DeepEqualNow(url, "")
}

func TestGetContext(t *testing.T) {
	a := assert.New(t)
	cli := New()

	baseCtx := context.Background()
	ctx, _ := cli.getContext(RequestOptions{
		Context: baseCtx,
	})
	a.DeepEqual(ctx, baseCtx)

	ctx, _ = cli.getContext(RequestOptions{})
	deadline, ok := ctx.Deadline()
	a.DeepEqualNow(ok, true)
	a.DeepEqualNow((deadline.UnixMilli() - time.Now().UnixMilli()), int64(1000))

	ctx, _ = cli.getContext(RequestOptions{
		Timeout: 3000,
	})
	deadline, ok = ctx.Deadline()
	a.DeepEqualNow(ok, true)
	a.DeepEqualNow((deadline.UnixMilli() - time.Now().UnixMilli()), int64(3000))

	ctx, _ = cli.getContext(RequestOptions{
		Timeout: RequestTimeoutNoLimit,
	})
	_, ok = ctx.Deadline()
	a.DeepEqualNow(ok, false)

	cli.timeout = 3000
	ctx, _ = cli.getContext(RequestOptions{})
	deadline, ok = ctx.Deadline()
	a.DeepEqualNow(ok, true)
	a.DeepEqualNow((deadline.UnixMilli() - time.Now().UnixMilli()), int64(3000))
}
