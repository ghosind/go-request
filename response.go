package request

import (
	"net/http"
	"strings"
)

func (cli *Client) decodeResponseBody(resp *http.Response, out any) error {
	contentType := resp.Header.Get("Content-Type")
	contentType, _, _ = strings.Cut(contentType, ";") // remove parameters
	contentType = strings.TrimSpace(contentType)

	switch contentType {
	case "application/json":
		return decodeJson(resp.Body, out)
	}

	return nil
}
