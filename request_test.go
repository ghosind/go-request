package request

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ghosind/go-assert"
)

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
