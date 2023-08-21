package request

import (
	"testing"

	"github.com/ghosind/go-assert"
)

func TestGetContentType(t *testing.T) {
	a := assert.New(t)

	a.Equal(getContentType(""), "json")
	a.Equal(getContentType("application/json"), "json")
	a.Equal(getContentType("application/json; charset=utf8"), "json")
	a.Equal(getContentType("application/xml; charset=utf8"), "unknown")
}
