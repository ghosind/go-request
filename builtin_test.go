package request

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/ghosind/go-assert"
	"github.com/ghosind/go-request/internal"
)

func TestMain(m *testing.M) {
	server := internal.NewMockServer()
	go server.Run()

	status := m.Run()

	server.Shutdown()
	os.Exit(status)
}

type testResponse struct {
	Path        *string `json:"path"`
	Method      *string `json:"method"`
	ContentType *string `json:"contentType"`
	Body        *string `json:"body"`
	Query       *string `json:"query"`
	Token       *string `json:"token"`
	UserAgent   *string `json:"userAgent"`
}

func TestRequestWithoutOptions(t *testing.T) {
	a := assert.New(t)

	data, _, err := ToString(Request("http://localhost:8080/test"))
	a.NilNow(err)

	a.NotDeepEqualNow(len(data), 0)

	payload := new(testResponse)
	a.NilNow(json.Unmarshal([]byte(data), &payload))

	a.NotNilNow(payload.Method)
	a.NotNilNow(payload.Path)
	a.DeepEqualNow(*payload.Method, "GET")
	a.DeepEqualNow(*payload.Path, "/test")
}

func TestRequestWithOptions(t *testing.T) {
	a := assert.New(t)

	data, _, err := ToObject[testResponse](Request("", RequestOptions{
		BaseURL: "http://localhost:8080",
	}))
	a.NilNow(err)

	a.NotNilNow(data.Method)
	a.DeepEqualNow(*data.Method, "GET")
}

func TestRequestMethods(t *testing.T) {
	a := assert.New(t)

	for method, fn := range map[string]func(string, ...RequestOptions) (*http.Response, error){
		"DELETE":  DELETE,
		"GET":     GET,
		"HEAD":    HEAD,
		"OPTIONS": OPTIONS,
		"PATCH":   PATCH,
		"POST":    POST,
		"PUT":     PUT,
	} {
		data, _, err := ToObject[testResponse](fn("http://localhost:8080"))
		a.NilNow(err, method)

		if method == "HEAD" {
			// no body for HEAD
			continue
		}

		a.NotNilNow(data.Method)
		a.DeepEqualNow(*data.Method, method)
	}
}

func TestRequestWithBody(t *testing.T) {
	a := assert.New(t)

	data, _, err := ToObject[testResponse](POST("http://localhost:8080", RequestOptions{
		Body: map[string]any{
			"data": "Hello world!",
		},
	}))
	a.NilNow(err)

	a.NotNilNow(data.Method)
	a.NotNilNow(data.ContentType)
	a.NotNilNow(data.Body)
	a.DeepEqualNow(*data.Method, "POST")
	a.DeepEqualNow(*data.ContentType, "application/json")
	a.DeepEqualNow(*data.Body, `{"data":"Hello world!"}`)
}

func TestRequestWithParameters(t *testing.T) {
	a := assert.New(t)

	data, _, err := ToObject[testResponse](GET("http://localhost:8080", RequestOptions{
		Parameters: map[string][]string{
			"q": {"test"},
			"v": {"1"},
		},
	}))
	a.NilNow(err)

	a.NotNilNow(data.Query)
	a.DeepEqualNow(*data.Query, "q=test&v=1")
}

func TestRequestWithHeader(t *testing.T) {
	a := assert.New(t)

	data, _, err := ToObject[testResponse](GET("http://localhost:8080", RequestOptions{
		Headers: map[string][]string{
			"Authorization": {"test-token"},
		},
	}))
	a.NilNow(err)

	a.NotNilNow(data.Token)
	a.DeepEqualNow(*data.Token, "test-token")

	cli := New(Config{
		Headers: map[string][]string{
			"Authorization": {"test-token"},
		},
	})
	data, _, err = ToObject[testResponse](cli.GET("http://localhost:8080"))
	a.NilNow(err)

	a.NotNilNow(data.Token)
	a.DeepEqualNow(*data.Token, "test-token")
}

func TestRequestWithUserAgent(t *testing.T) {
	a := assert.New(t)

	data, _, err := ToObject[testResponse](GET("http://localhost:8080"))
	a.NilNow(err)

	a.NotNilNow(data.UserAgent)
	a.DeepEqualNow(*data.UserAgent, RequestDefaultUserAgent)

	cli := New(Config{
		UserAgent: "Test-client",
	})
	data, _, err = ToObject[testResponse](cli.GET("http://localhost:8080"))
	a.NilNow(err)

	a.NotNilNow(data.UserAgent)
	a.DeepEqualNow(*data.UserAgent, "Test-client")

	data, _, err = ToObject[testResponse](cli.GET("http://localhost:8080", RequestOptions{
		UserAgent: "Another-test-client",
	}))
	a.NilNow(err)

	a.NotNilNow(data.UserAgent)
	a.DeepEqualNow(*data.UserAgent, "Another-test-client")
}

func TestRequestWithInvalidConfig(t *testing.T) {
	a := assert.New(t)

	_, err := Request("http://example.com", RequestOptions{
		Method: "UNKNOWN",
	})
	a.NotNilNow(err)
}
