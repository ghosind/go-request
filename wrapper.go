package request

import (
	"encoding/json"
	"io"
	"net/http"
)

func ToObject[T any](resp *http.Response, err error) (*T, error) {
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	out := new(T)

	contentType := resp.Header.Get("Content-Type")
	switch contentType {
	case "application/json":
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			return nil, err
		}
	}

	return out, nil
}

func ToString(resp *http.Response, err error) (string, error) {
	if err != nil {
		return "", err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	resp.Body.Close()

	return string(data), nil
}
