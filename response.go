package request

import (
	"net/http"
)

func (cli *Client) decodeResponseBody(resp *http.Response, out any) error {
	contentType := getContentType(resp.Header.Get("Content-Type"))

	switch contentType {
	case "application/json":
		return decodeJson(resp.Body, out)
	}

	return nil
}
