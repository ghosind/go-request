package request

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/ghosind/go-assert"
)

func TestToObject(t *testing.T) {
	a := assert.New(t)

	_, _, err := ToObject[testResponse](nil, errors.New("test error"))
	a.NotNilNow(err)

	_, _, err = ToObject[testResponse](nil, nil)
	a.NotNilNow(err)

	_, _, err = ToObject[testResponse](&http.Response{}, nil)
	a.NotNilNow(err)

	_, _, err = ToObject[testResponse](&http.Response{
		Body: io.NopCloser(bytes.NewReader([]byte("Not a JSON"))),
	}, nil)
	a.NotNilNow(err)

	data, _, err := ToObject[testResponse](&http.Response{
		Body: io.NopCloser(bytes.NewReader([]byte(""))),
	}, nil)
	a.NilNow(err)
	a.NilNow(data)

	data, _, err = ToObject[testResponse](&http.Response{
		Body: io.NopCloser(bytes.NewReader([]byte(`{"method":"GET"}`))),
	}, nil)
	a.NilNow(err)
	a.NotNilNow(data.Method)
	a.DeepEqualNow(*data.Method, "GET")
}

func TestToString(t *testing.T) {
	a := assert.New(t)

	_, _, err := ToString(nil, errors.New("test error"))
	a.NotNilNow(err)

	_, _, err = ToString(nil, nil)
	a.NotNilNow(err)

	_, _, err = ToString(&http.Response{}, nil)
	a.NotNilNow(err)

	data, _, err := ToString(&http.Response{
		Body: io.NopCloser(bytes.NewReader([]byte("Hello world!"))),
	}, nil)
	a.NilNow(err)
	a.DeepEqualNow(data, "Hello world!")
}
