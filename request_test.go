package request

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/ghosind/go-assert"
)

type ExampleProduct struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Category  string `json:"category"`
	IsDeleted *bool  `json:"isDeleted,omitempty"`
}

type ExampleProductsSearchData struct {
	Products *[]ExampleProduct
}

type AuthData struct {
	Token string `json:"token"`
}

func TestSimpleRequest(t *testing.T) {
	a := assert.New(t)

	data, _, err := ToString(Request("https://dummyjson.com/products/1"))
	if e := a.Nil(err); e != nil {
		a.FailNow()
	}

	if e := a.NotDeepEqual(len(data), 0); e != nil {
		a.FailNow()
	}

	product := new(ExampleProduct)
	if e := a.Nil(json.Unmarshal([]byte(data), &product)); e != nil {
		a.FailNow()
	}

	a.DeepEqual(product.ID, int64(1))
}

func TestGetRequest(t *testing.T) {
	a := assert.New(t)

	data, _, err := ToObject[ExampleProduct](GET("/products/1", RequestOptions{
		BaseURL: "https://dummyjson.com",
	}))
	if e := a.Nil(err); e != nil {
		a.FailNow()
	}

	a.DeepEqual(data.ID, int64(1))
}

func TestPOSTRequest(t *testing.T) {
	a := assert.New(t)
	title := "MacBook Pro"

	data, _, err := ToObject[ExampleProduct](POST("/products/add", RequestOptions{
		BaseURL: "https://dummyjson.com/",
		Timeout: RequestTimeoutNoLimit,
		Body: map[string]any{
			"title": title,
		},
	}))
	if e := a.Nil(err); e != nil {
		a.FailNow()
	}

	a.DeepEqual(data.Title, title)
}

func TestPUTRequest(t *testing.T) {
	a := assert.New(t)
	title := "Apple"

	data, _, err := ToObject[ExampleProduct](PUT("/products/1", RequestOptions{
		BaseURL: "https://dummyjson.com/",
		Timeout: RequestTimeoutNoLimit,
		Body: map[string]any{
			"title": title,
		},
	}))
	if e := a.Nil(err); e != nil {
		a.FailNow()
	}

	a.DeepEqual(data.Title, title)
}

func TestDeleteRequest(t *testing.T) {
	a := assert.New(t)
	data, _, err := ToObject[ExampleProduct](DELETE("/products/1", RequestOptions{
		BaseURL: "https://dummyjson.com/",
		Timeout: RequestTimeoutNoLimit,
	}))
	if e := a.Nil(err); e != nil {
		a.FailNow()
	}

	if e := a.NotNil(data.IsDeleted); e == nil {
		a.DeepEqual(*data.IsDeleted, true)
	}
}

func TestGETRequestWithParameters(t *testing.T) {
	a := assert.New(t)

	data, _, err := ToObject[ExampleProductsSearchData](GET("/products/search", RequestOptions{
		BaseURL: "https://dummyjson.com",
		Parameters: map[string][]string{
			"q": {"laptop"},
		},
		Timeout: 3 * 1000,
	}))
	if e := a.Nil(err); e != nil {
		a.FailNow()
	}

	if e := a.NotNil(data.Products); e != nil {
		a.FailNow()
	}

	for _, product := range *data.Products {
		a.DeepEqual(strings.Contains(product.Category, "laptop"), true)
	}
}

func TestRequestWithHeader(t *testing.T) {
	a := assert.New(t)
	cli := New(Config{
		BaseURL: "https://dummyjson.com",
	})

	authData, _, err := ToObject[AuthData](cli.POST("/auth/login", RequestOptions{
		Body: map[string]any{
			"username": "atuny0",
			"password": "9uQFF1Lh",
		},
	}))
	if e := a.Nil(err); e != nil {
		a.FailNow()
	}

	cli.Headers["Authorization"] = []string{authData.Token}

	data, resp, err := ToObject[ExampleProduct](cli.GET("/auth/products/1"))
	if e := a.Nil(err); e != nil {
		a.FailNow()
	}

	a.DeepEqual(resp.StatusCode, http.StatusOK)
	a.DeepEqual(data.ID, int64(1))
}
