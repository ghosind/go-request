package request

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
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
	data, _, err := ToString(Request("https://dummyjson.com/products/1"))
	if err != nil {
		t.Fatalf("Unexpected request error: %v", err)
	}

	if len(data) == 0 {
		t.Fatalf("Expect data's length is not 0")
	}

	product := new(ExampleProduct)
	if err := json.Unmarshal([]byte(data), &product); err != nil {
		t.Fatalf("Unexpected unmarshal error: %v", err)
	}

	if product.ID != 1 {
		t.Fatalf("Expect product's id is 1, actually %d", product.ID)
	}
}

func TestGetRequest(t *testing.T) {
	data, _, err := ToObject[ExampleProduct](GET("/products/1", RequestOptions{
		BaseURL: "https://dummyjson.com",
	}))
	if err != nil {
		t.Fatalf("Unexpected request error: %v", err)
	}

	if data.ID != 1 {
		t.Fatalf("Expect product's id is 1, actually %d", data.ID)
	}
}

func TestPOSTRequest(t *testing.T) {
	title := "MacBook Pro"

	data, _, err := ToObject[ExampleProduct](POST("/products/add", RequestOptions{
		BaseURL: "https://dummyjson.com/",
		Timeout: RequestTimeoutNone,
		Body: map[string]any{
			"title": title,
		},
	}))
	if err != nil {
		t.Fatalf("Unexpected request error: %v", err)
	}

	if data.Title != title {
		t.Fatalf("Expect product's title is \"%s\", actually \"%s\"", title, data.Title)
	}
}

func TestPUTRequest(t *testing.T) {
	title := "Apple"

	data, _, err := ToObject[ExampleProduct](PUT("/products/1", RequestOptions{
		BaseURL: "https://dummyjson.com/",
		Timeout: RequestTimeoutNone,
		Body: map[string]any{
			"title": title,
		},
	}))
	if err != nil {
		t.Fatalf("Unexpected request error: %v", err)
	}

	if data.Title != title {
		t.Fatalf("Expect product's title is \"%s\", actually \"%s\"", title, data.Title)
	}
}

func TestDeleteRequest(t *testing.T) {
	data, _, err := ToObject[ExampleProduct](DELETE("/products/1", RequestOptions{
		BaseURL: "https://dummyjson.com/",
		Timeout: RequestTimeoutNone,
	}))
	if err != nil {
		t.Fatalf("Unexpected request error: %v", err)
	}

	if data.IsDeleted == nil || *data.IsDeleted != true {
		t.Fatalf("Expect product's isDelete is true")
	}
}

func TestGETRequestWithParameters(t *testing.T) {
	data, _, err := ToObject[ExampleProductsSearchData](GET("/products/search", RequestOptions{
		BaseURL: "https://dummyjson.com",
		Parameters: map[string][]string{
			"q": {"laptop"},
		},
		Timeout: 3 * 1000,
	}))
	if err != nil {
		t.Fatalf("Unexpected request error: %v", err)
	}

	if data.Products != nil {
		for _, product := range *data.Products {
			if !strings.Contains(product.Category, "laptop") {
				t.Fatalf("Expect product's category contains 'laptop', actually %s", product.Category)
			}
		}
	}
}

func TestRequestWithHeader(t *testing.T) {
	cli := New(Config{
		BaseURL: "https://dummyjson.com",
	})

	authData, _, err := ToObject[AuthData](cli.POST("/auth/login", RequestOptions{
		Body: map[string]any{
			"username": "atuny0",
			"password": "9uQFF1Lh",
		},
	}))
	if err != nil {
		t.Fatalf("Unexpected request error: %v", err)
	}

	cli.Headers["Authorization"] = []string{authData.Token}

	data, resp, err := ToObject[ExampleProduct](cli.GET("/auth/products/1"))
	if err != nil {
		t.Fatalf("Unexpected request error: %v", err)
	} else if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected HTTP status code: %d(%s)", resp.StatusCode, resp.Status)
	}

	if data.ID != 1 {
		t.Fatalf("Expect product's id is 1, actually %d", data.ID)
	}
}
