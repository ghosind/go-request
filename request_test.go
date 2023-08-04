package request

import (
	"net/http"
	"strings"
	"testing"

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
