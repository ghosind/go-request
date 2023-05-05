package request

import (
	"encoding/json"
	"io"
	"net/http"
)

func ToObject[T any](resp *http.Response, err error) (*T, *http.Response, error) {
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	out := new(T)

	contentType := resp.Header.Get("Content-Type")
	switch getContentType(contentType) {
	case "application/json":
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			return nil, resp, err
		}
	}

	return out, resp, nil
}

func ToString(resp *http.Response, err error) (string, *http.Response, error) {
	if err != nil {
		return "", nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	resp.Body.Close()

	return string(data), resp, nil
}
