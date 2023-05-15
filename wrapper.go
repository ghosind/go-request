package request

import (
	"encoding/json"
	"io"
	"net/http"
)

// ToObject reads data from the response body and tries to decode it to an object as the parameter
// type. It'll read the encoding type from the 'Content-Type' field in the response header. The
// method will close the body of the response that after read.
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

// ToString reads data from the response body and returns them as a string. The method will close
// the body of the response that after read.
func ToString(resp *http.Response, err error) (string, *http.Response, error) {
	if err != nil {
		return "", nil, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", resp, err
	}

	return string(data), resp, nil
}
