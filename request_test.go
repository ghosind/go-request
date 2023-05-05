package request

import (
	"net/http"
	"strings"
	"testing"
)

type ExampleProduct struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category"`
}

type ExampleProductsSearchData struct {
	Products *[]ExampleProduct
}

func TestSimpleRequest(t *testing.T) {
	data, _, err := ToObject[ExampleProduct](Request("https://dummyjson.com/products/1", RequestOptions{
		Timeout: 3000,
	}))
	if err != nil {
		t.Fatalf("Unexpected request error: %v", err)
	}

	if data.ID != 1 {
		t.Fatalf("Expect product id is 1, actually %d", data.ID)
	}
}

func TestSimplePOSTRequest(t *testing.T) {
	title := "MacBook Pro"

	data, _, err := ToObject[ExampleProduct](Request("/products/add", RequestOptions{
		BaseURL: "https://dummyjson.com/",
		Method:  http.MethodPost,
		Timeout: 3000,
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

func TestRequestWithParameters(t *testing.T) {
	data, _, err := ToObject[ExampleProductsSearchData](Request("/products/search", RequestOptions{
		BaseURL: "https://dummyjson.com",
		Parameters: map[string][]string{
			"q": {"laptop"},
		},
	}))
	if err != nil {
		t.Fatalf("Unexpected request error: %v", err)
	}

	if data.Products != nil {
		for _, product := range *data.Products {
			if !strings.Contains(product.Category, "laptop") {
				t.Fatalf("Expect product category contains 'laptop', actually %s", product.Category)
			}
		}
	}
}
