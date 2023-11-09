package request

import (
	"errors"
	"net/http"
	"testing"

	"github.com/ghosind/go-assert"
)

func TestAddInterceptors(t *testing.T) {
	a := assert.New(t)

	reqInterIds := UseRequestInterceptor(func(r *http.Request) error {
		// just a test interceptor
		return nil
	}, nil, func(r *http.Request) error {
		// just a test interceptor
		return nil
	})
	respInterIds := UseResponseInterceptor(func(r *http.Response) error {
		// just a test interceptor
		return nil
	}, nil, func(r *http.Response) error {
		// just a test interceptor
		return nil
	})
	a.TrueNow(len(reqInterIds) == 3)
	a.TrueNow(len(respInterIds) == 3)
	a.EqualNow(len(defaultClient.reqInterceptors), 2)
	a.EqualNow(len(defaultClient.respInterceptors), 2)

	for _, id := range reqInterIds {
		ok := RemoveRequestInterceptor(id)
		a.EqualNow(ok, id != 0)
	}
	for _, id := range respInterIds {
		ok := RemoveResponseInterceptor(id)
		a.EqualNow(ok, id != 0)
	}

	// Remove again, should be all fail
	for _, id := range reqInterIds {
		ok := RemoveRequestInterceptor(id)
		a.NotTrueNow(ok)
	}
	for _, id := range respInterIds {
		ok := RemoveResponseInterceptor(id)
		a.NotTrueNow(ok)
	}
}

func TestRequestWithInterceptors(t *testing.T) {
	a := assert.New(t)

	reqIntercepted := false
	respIntercepted := false

	reqInterId := UseRequestInterceptor(func(r *http.Request) error {
		reqIntercepted = true
		return nil
	})
	respInterId := UseResponseInterceptor(func(r *http.Response) error {
		respIntercepted = true
		return nil
	})
	a.TrueNow(len(reqInterId) == 1)
	a.TrueNow(len(respInterId) == 1)

	_, err := Request("http://localhost:8080")
	a.NilNow(err)
	a.TrueNow(reqIntercepted)
	a.TrueNow(respIntercepted)

	ok := RemoveRequestInterceptor(reqInterId[0])
	a.TrueNow(ok)
	ok = RemoveResponseInterceptor(respInterId[0])
	a.TrueNow(ok)
}

func TestRequestInterceptorFailure(t *testing.T) {
	a := assert.New(t)

	reqIntercepted := false

	reqInterIds := UseRequestInterceptor(func(r *http.Request) error {
		return errors.New("expected error")
	}, func(r *http.Request) error {
		reqIntercepted = true
		return nil
	})
	a.TrueNow(len(reqInterIds) == 2)

	_, err := Request("http://localhost:8080")
	a.NotNilNow(err)
	a.EqualNow(err.Error(), "expected error")
	a.NotTrueNow(reqIntercepted)

	for _, id := range reqInterIds {
		RemoveRequestInterceptor(id)
	}
}

func TestResponseInterceptorFailure(t *testing.T) {
	a := assert.New(t)

	respIntercepted := false

	respInterIds := UseResponseInterceptor(func(r *http.Response) error {
		return errors.New("expected error")
	}, func(r *http.Response) error {
		respIntercepted = true
		return nil
	})
	a.TrueNow(len(respInterIds) == 2)

	resp, err := Request("http://localhost:8080")
	a.NotNilNow(err)
	a.EqualNow(err.Error(), "expected error")
	a.NotNilNow(resp)
	a.NotTrueNow(respIntercepted)

	for _, id := range respInterIds {
		RemoveResponseInterceptor(id)
	}
}
