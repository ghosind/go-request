package request

import (
	"testing"

	"github.com/ghosind/go-assert"
)

func TestGetRequestBody(t *testing.T) {
	a := assert.New(t)
	cli := New()

	out, err := cli.getRequestBody(RequestOptions{})
	a.NilNow(err)
	a.NilNow(out)

	_, err = cli.getRequestBody(RequestOptions{
		Body:        []string{"Test"},
		ContentType: "Unknown",
	})
	a.NotNilNow(err)
}

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

	out, err := cli.encodeRequestBody(nil, "")
	a.NilNow(err)
	a.NilNow(out)

	_, err = cli.encodeRequestBody(testStruct{
		Message: "Hello",
	}, "unknown")
	a.NotNilNow(err)
}

func testEncodeRequestBody(a *assert.Assertion, cli *Client, data any, expect []byte) {
	out, err := cli.encodeRequestBody(data, "")
	a.NilNow(err)
	a.DeepEqual(out, expect)
}
