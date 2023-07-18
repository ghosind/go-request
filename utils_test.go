package request

import (
	"testing"

	"github.com/ghosind/go-assert"
)

func TestGetContentType(t *testing.T) {
	a := assert.New(t)

	a.DeepEqual(getContentType(""), "")
	a.DeepEqual(getContentType("application/json"), "application/json")
	a.DeepEqual(getContentType("application/json; charset=utf8"), "application/json")
}
