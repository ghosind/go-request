package request

import (
	"reflect"
	"testing"
)

func TestEncodeRequestBody(t *testing.T) {
	cli := New()

	testEncodeRequestBody(t, cli, "Test", []byte("Test"))
	testEncodeRequestBody(t, cli, []byte("Test"), []byte("Test"))
	testEncodeRequestBody(t, cli, map[string]any{
		"message": "Hello",
	}, []byte(`{"message":"Hello"}`))

	type testStruct struct {
		Message string `json:"message"`
	}
	testEncodeRequestBody(t, cli, testStruct{
		Message: "Hello",
	}, []byte(`{"message":"Hello"}`))
}

func testEncodeRequestBody(t *testing.T, cli *Client, data any, expect []byte) {
	out, err := cli.encodeRequestBody(data, "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(out, expect) {
		t.Errorf("encodeRequestBody returns \"%s\", expect \"%s\"", string(out), string(expect))
	}
}
