package request

import (
	"testing"

	"github.com/ghosind/go-assert"
)

func TestEncodeRequestBody(t *testing.T) {
	a := assert.New(t)
	cli := New()

	testEncodeRequestBody(a, cli, "Test", []byte("Test"))
	testEncodeRequestBody(a, cli, []byte("Test"), []byte("Test"))
	testEncodeRequestBody(a, cli, map[string]any{
		"message": "Hello",
	}, []byte(`{"message":"Hello"}`))

	type testStruct struct {
		Message string `json:"message"`
	}
	testEncodeRequestBody(a, cli, testStruct{
		Message: "Hello",
	}, []byte(`{"message":"Hello"}`))
}

func testEncodeRequestBody(a *assert.Assertion, cli *Client, data any, expect []byte) {
	out, err := cli.encodeRequestBody(data, "")
	if err != nil {
		a.Errorf("encodeRequestBody()'s error = %v, want nil", err)
	} else {
		a.DeepEqual(out, expect)
	}
}
