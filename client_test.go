package request

import (
	"testing"

	"github.com/ghosind/go-assert"
)

func TestCreateClientWithoutNewFunction(t *testing.T) {
	a := assert.New(t)

	cli := &Client{}

	_, err := cli.GET("http://localhost:8080")
	a.NilNow(err)
}
