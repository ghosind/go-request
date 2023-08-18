package request

import (
	"encoding/json"
	"io"
	"net/http"
)

// ToObject reads data from the response body and tries to decode it to an object as the parameter
// type. It'll read the encoding type from the 'Content-Type' field in the response header. The
// method will close the body of the response that after read.
//
//	data, resp, err := request.ToObject[SomeStructType](request.Request("https://example.com"))
//	if err != nil {
//	  // Error handling
//	}
//	// Data or response handling
func ToObject[T any](resp *http.Response, err error) (*T, *http.Response, error) {
	if err != nil {
		return nil, nil, err
	}
	if resp == nil || resp.Body == nil {
		return nil, nil, ErrInvalidResp
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, err
	} else if len(data) == 0 {
		return nil, resp, nil
	}

	out := new(T)

	contentType := resp.Header.Get("Content-Type")
	switch getContentType(contentType) {
	default:
		if err := json.Unmarshal(data, &out); err != nil {
			return nil, resp, err
		}
	}

	return out, resp, nil
}

// ToString reads data from the response body and returns them as a string. The method will close
// the body of the response that after read.
//
//	content, resp, err := request.ToString(request.Request("https://example.com"))
//	if err != nil {
//	  // Error handling
//	}
//	// Response handling
func ToString(resp *http.Response, err error) (string, *http.Response, error) {
	if err != nil {
		return "", nil, err
	}
	if resp == nil || resp.Body == nil {
		return "", nil, ErrInvalidResp
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", resp, err
	}

	return string(data), resp, nil
}
