package request

import (
	"encoding/json"
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
}

func TestSimpleRequest(t *testing.T) {
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

func TestGetRequest(t *testing.T) {
	a := assert.New(t)

	data, _, err := ToObject[testResponse](GET("", RequestOptions{
		BaseURL: "http://localhost:8080",
	}))
	a.NilNow(err)

	a.NotNilNow(data.Method)
	a.DeepEqualNow(*data.Method, "GET")
}

func TestPOSTRequest(t *testing.T) {
	a := assert.New(t)

	data, _, err := ToObject[testResponse](POST("", RequestOptions{
		BaseURL: "http://localhost:8080",
		Timeout: RequestTimeoutNoLimit,
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

func TestPUTRequest(t *testing.T) {
	a := assert.New(t)

	data, _, err := ToObject[testResponse](PUT("/", RequestOptions{
		BaseURL: "http://localhost:8080",
		Timeout: RequestTimeoutNoLimit,
		Body: map[string]any{
			"data": "Hello world!",
		},
	}))
	a.NilNow(err)

	a.NotNilNow(data.Method)
	a.NotNilNow(data.Body)
	a.DeepEqualNow(*data.Method, "PUT")
	a.DeepEqualNow(*data.Body, `{"data":"Hello world!"}`)
}

func TestDeleteRequest(t *testing.T) {
	a := assert.New(t)
	data, _, err := ToObject[testResponse](DELETE("/", RequestOptions{
		BaseURL: "http://localhost:8080",
		Timeout: RequestTimeoutNoLimit,
	}))
	a.NilNow(err)

	a.NotNilNow(data.Method)
	a.DeepEqualNow(*data.Method, "DELETE")
}

func TestGETRequestWithParameters(t *testing.T) {
	a := assert.New(t)

	data, _, err := ToObject[testResponse](GET("", RequestOptions{
		BaseURL: "http://localhost:8080",
		Parameters: map[string][]string{
			"q": {"test"},
		},
		Timeout: 3 * 1000,
	}))
	a.NilNow(err)

	a.NotNilNow(data.Query)
	a.DeepEqualNow(*data.Query, "q=test")
}

func TestRequestWithHeader(t *testing.T) {
	a := assert.New(t)
	cli := New(Config{
		BaseURL: "http://localhost:8080",
	})

	cli.Headers["Authorization"] = []string{"test-token"}

	data, _, err := ToObject[testResponse](cli.GET("/"))
	a.NilNow(err)

	a.NotNilNow(data.Token)
	a.DeepEqualNow(*data.Token, "test-token")
}
