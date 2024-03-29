package request

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/ghosind/go-assert"
	"github.com/ghosind/go-request/internal"
)

func TestMain(m *testing.M) {
	server := internal.NewMockServer()
	go server.Run()

	proxy := internal.NewProxyServer()
	go proxy.Run()

	status := m.Run()

	server.Shutdown()
	os.Exit(status)
}

type testResponse struct {
	Path        *string              `json:"path"`
	Method      *string              `json:"method"`
	ContentType *string              `json:"contentType"`
	Body        *string              `json:"body"`
	Query       *string              `json:"query"`
	Token       *string              `json:"token"`
	UserAgent   *string              `json:"userAgent"`
	Headers     *map[string][]string `json:"headers"`
}

func customParametersSerializer(params map[string][]string) string {
	sb := strings.Builder{}

	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		if sb.Len() != 0 {
			sb.WriteRune('&')
		}

		v := params[k]

		if len(v) == 1 {
			sb.WriteString(url.QueryEscape(k))
			sb.WriteRune('=')
			sb.WriteString(url.QueryEscape(v[0]))
		} else {
			for i, vv := range v {
				if i > 0 {
					sb.WriteRune('&')
				}
				sb.WriteString(url.QueryEscape(k))
				sb.WriteRune('[')
				sb.WriteString(strconv.Itoa(i))
				sb.WriteRune(']')
				sb.WriteRune('=')
				sb.WriteString(url.QueryEscape(vv))
			}
		}
	}

	return sb.String()
}

func TestRequestWithoutOptions(t *testing.T) {
	a := assert.New(t)

	content, _, err := ToString(Request("http://localhost:8080/test"))
	a.NilNow(err)

	a.NotEqualNow(len(content), 0)

	data := new(testResponse)
	a.NilNow(json.Unmarshal([]byte(content), &data))

	a.NotNilNow(data.Method)
	a.NotNilNow(data.Path)
	a.EqualNow(*data.Method, "GET")
	a.EqualNow(*data.Path, "/test")
}

func TestRequestWithOptions(t *testing.T) {
	a := assert.New(t)

	data, _, err := ToObject[testResponse](Request("/test", RequestOptions{
		BaseURL: "http://localhost:8080",
	}))
	a.NilNow(err)

	a.NotNilNow(data.Method)
	a.NotNilNow(data.Path)
	a.EqualNow(*data.Method, "GET")
	a.EqualNow(*data.Path, "/test")
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
		a.EqualNow(*data.Method, method)
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
	a.EqualNow(*data.Method, "POST")
	a.EqualNow(*data.ContentType, "application/json")
	a.EqualNow(*data.Body, `{"data":"Hello world!"}`)
}

func TestRequestWithContext(t *testing.T) {
	a := assert.New(t)

	_, err := GET("http://localhost:8080", RequestOptions{
		Context: context.Background(),
	})
	a.NilNow(err)

	ctx, canFunc := context.WithCancel(context.Background())
	canFunc()
	_, err = GET("http://localhost:8080", RequestOptions{
		Context:     ctx,
		ContentType: "unknown",
	})
	a.NotNilNow(err)
}

func TestRequestWithGZipEncodedBody(t *testing.T) {
	a := assert.New(t)

	data, resp, err := ToObject[testResponse](Request("", RequestOptions{
		BaseURL: "http://localhost:8080",
		Headers: map[string][]string{
			"Accept-Encoding": {"gzip", "deflate"},
		},
	}))
	a.NilNow(err)

	a.NotNilNow(data.Method)
	a.EqualNow(*data.Method, "GET")
	a.NotTrueNow(resp.Header.Get("Content-Encoding"))
}

func TestRequestWithDeflateEncodedBody(t *testing.T) {
	a := assert.New(t)

	data, resp, err := ToObject[testResponse](Request("", RequestOptions{
		BaseURL: "http://localhost:8080",
		Headers: map[string][]string{
			"Accept-Encoding": {"deflate"},
		},
	}))
	a.NilNow(err)

	a.NotNilNow(data.Method)
	a.EqualNow(*data.Method, "GET")
	a.NotTrueNow(resp.Header.Get("Content-Encoding"))
}

func TestRequestWithInvalidContentEncoding(t *testing.T) {
	a := assert.New(t)

	_, _, err := ToObject[testResponse](Request("", RequestOptions{
		BaseURL: "http://localhost:8080",
		Headers: map[string][]string{
			"Accept-Encoding": {"gzip"},
		},
		Parameters: map[string][]string{
			"contentEncoding": {"deflate"}, // force invalid encoding
		},
	}))
	a.NotNilNow(err)

	_, _, err = ToString(Request("", RequestOptions{
		BaseURL: "http://localhost:8080",
		Headers: map[string][]string{
			"Accept-Encoding": {"gzip"},
		},
		Parameters: map[string][]string{
			"contentEncoding": {"deflate"}, // force invalid encoding
		},
	}))
	a.NotNilNow(err)
}

func TestRequestWithContentEncodingAndDisableDecompress(t *testing.T) {
	a := assert.New(t)

	resp, err := Request("", RequestOptions{
		BaseURL: "http://localhost:8080",
		Headers: map[string][]string{
			"Accept-Encoding": {"gzip", "deflate"},
		},
		DisableDecompress: true,
	})
	a.NilNow(err)

	a.EqualNow(resp.Header.Get("Content-Encoding"), "gzip")
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
	a.EqualNow(*data.Query, "q=test&v=1")
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
	a.EqualNow(*data.Token, "test-token")

	cli := New(Config{
		Headers: map[string][]string{
			"Authorization": {"test-token"},
		},
	})
	data, _, err = ToObject[testResponse](cli.GET("http://localhost:8080"))
	a.NilNow(err)

	a.NotNilNow(data.Token)
	a.EqualNow(*data.Token, "test-token")
}

func TestRequestWithUserAgent(t *testing.T) {
	a := assert.New(t)

	data, _, err := ToObject[testResponse](GET("http://localhost:8080"))
	a.NilNow(err)

	a.NotNilNow(data.UserAgent)
	a.EqualNow(*data.UserAgent, RequestDefaultUserAgent)

	cli := New(Config{
		UserAgent: "Test-client",
	})
	data, _, err = ToObject[testResponse](cli.GET("http://localhost:8080"))
	a.NilNow(err)

	a.NotNilNow(data.UserAgent)
	a.EqualNow(*data.UserAgent, "Test-client")

	data, _, err = ToObject[testResponse](cli.GET("http://localhost:8080", RequestOptions{
		UserAgent: "Another-test-client",
	}))
	a.NilNow(err)

	a.NotNilNow(data.UserAgent)
	a.EqualNow(*data.UserAgent, "Another-test-client")
}

func TestRequestWithInvalidConfig(t *testing.T) {
	a := assert.New(t)

	_, err := Request("http://example.com", RequestOptions{
		Method: "UNKNOWN",
	})
	a.NotNilNow(err)
}

func TestRequestWithMaxRedirects(t *testing.T) {
	a := assert.New(t)

	resp, err := Request("http://localhost:8080/redirect", RequestOptions{
		MaxRedirects: 3,
	})
	a.NilNow(err)

	locationUrl := resp.Header.Get("Location")
	location, err := url.Parse(locationUrl)
	a.NilNow(err)

	tried := location.Query().Get("tried")
	a.EqualNow(tried, "3")
}

func TestRequestWithDefaultMaxRedirects(t *testing.T) {
	a := assert.New(t)

	resp, err := Request("http://localhost:8080/redirect")
	a.NilNow(err)

	locationUrl := resp.Header.Get("Location")
	location, err := url.Parse(locationUrl)
	a.NilNow(err)

	tried := location.Query().Get("tried")
	a.EqualNow(tried, strconv.FormatInt(int64(RequestDefaultMaxRedirects), 10))
}

func TestRequestWithClientMaxRedirects(t *testing.T) {
	a := assert.New(t)

	cli := New(Config{
		MaxRedirects: 3,
	})

	resp, err := cli.Request("http://localhost:8080/redirect")
	a.NilNow(err)

	locationUrl := resp.Header.Get("Location")
	location, err := url.Parse(locationUrl)
	a.NilNow(err)

	tried := location.Query().Get("tried")
	a.EqualNow(tried, "3")
}

func TestRequestWithNoRedirects(t *testing.T) {
	a := assert.New(t)

	resp, err := Request("http://localhost:8080/redirect", RequestOptions{
		MaxRedirects: RequestNoRedirects,
	})
	a.NilNow(err)

	locationUrl := resp.Header.Get("Location")
	location, err := url.Parse(locationUrl)
	a.NilNow(err)

	tried := location.Query().Get("tried")
	a.EqualNow(tried, "1")
}

func TestRequestWithParametersSerializer(t *testing.T) {
	a := assert.New(t)

	data, _, err := ToObject[testResponse](Request("http://localhost:8080", RequestOptions{
		Parameters: map[string][]string{
			"status": {"0", "10"},
		},
		ParametersSerializer: customParametersSerializer,
	}))
	a.NilNow(err)
	a.NotNilNow(data)
	a.NotNilNow(data.Query)
	a.Equal(*data.Query, "status[0]=0&status[1]=10")

	data, _, err = ToObject[testResponse](Request("http://localhost:8080", RequestOptions{
		Parameters: map[string][]string{
			"status": {"0", "10"},
		},
	}))
	a.NilNow(err)
	a.NotNilNow(data)
	a.NotNilNow(data.Query)
	a.Equal(*data.Query, "status=0&status=10")

	cli := New(Config{
		ParametersSerializer: customParametersSerializer,
	})
	data, _, err = ToObject[testResponse](cli.Request("http://localhost:8080", RequestOptions{
		Parameters: map[string][]string{
			"status": {"0", "10"},
		},
	}))
	a.NilNow(err)
	a.NotNilNow(data)
	a.NotNilNow(data.Query)
	a.Equal(*data.Query, "status[0]=0&status[1]=10")
}

func TestRequestWithDefaultValidateStatus(t *testing.T) {
	a := assert.New(t)

	resp, err := Request("http://localhost:8080/status", RequestOptions{
		Parameters: map[string][]string{
			"status": {"400"},
		},
	})
	a.NotNilNow(err)
	a.NotNilNow(resp)
	a.Equal(resp.StatusCode, 400)
}

func TestRequestWithCustomValidateStatus(t *testing.T) {
	a := assert.New(t)

	resp, err := Request("http://localhost:8080/status", RequestOptions{
		Parameters: map[string][]string{
			"status": {"400"},
		},
		ValidateStatus: func(status int) bool {
			return status == 400
		},
	})
	a.NilNow(err)
	a.NotNilNow(resp)
	a.Equal(resp.StatusCode, 400)
}

func TestRequestWithClientValidateStatus(t *testing.T) {
	a := assert.New(t)

	cli := New(Config{
		ValidateStatus: func(status int) bool {
			return status == 400
		},
	})

	resp, err := cli.Request("http://localhost:8080/status", RequestOptions{
		Parameters: map[string][]string{
			"status": {"400"},
		},
	})
	a.NilNow(err)
	a.NotNilNow(resp)
	a.Equal(resp.StatusCode, 400)
}
