package request

import (
	"strings"
	"testing"
)

type ExampleProduct struct {
	Id       int64
	Category string
}

type ExampleProductsSearchData struct {
	Products *[]ExampleProduct
}

func TestSimpleRequest(t *testing.T) {
	var data ExampleProduct

	err := Request("https://dummyjson.com/products/1", &data)
	if err != nil {
		t.Fatalf("Unexpected request error: %v", err)
	}

	if data.Id != 1 {
		t.Fatalf("Expect product id is 1, actually %d", data.Id)
	}
}

func TestRequestWithParameters(t *testing.T) {
	var data ExampleProductsSearchData

	err := Request("https://dummyjson.com/products/search", &data, RequestOptions{
		Parameters: map[string][]string{
			"q": {"laptop"},
		},
	})
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
