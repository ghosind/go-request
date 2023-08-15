package request

import (
	"testing"

	"github.com/ghosind/go-assert"
)

func TestGetContentType(t *testing.T) {
	a := assert.New(t)

	a.Equal(getContentType(""), "")
	a.Equal(getContentType("application/json"), "application/json")
	a.Equal(getContentType("application/json; charset=utf8"), "application/json")
}
