package request

import (
	"bytes"
	"encoding/json"
	"io"
)

// encodeRequestBody encodes the request body by the specific encoder. It'll get the encoding
// type from the 'ContentType' field in the request config, and it uses JSON as the default
// encoder. It will skip encoding the request body data if it is a byte array slice or a string.
func (cli *Client) encodeRequestBody(opt RequestOptions) (io.Reader, error) {
	body := opt.Body

	if body == nil {
		return nil, nil
	}

	// Encoding byte slice or string is unnecessary.
	if v, ok := opt.Body.([]byte); ok {
		return bytes.NewBuffer(v), nil
	} else if v, ok := opt.Body.(string); ok {
		return bytes.NewBuffer([]byte(v)), nil
	}

	switch getContentType(opt.ContentType) {
	default:
		return encodeBodyToJSON(body)
	}
}

// encodeJSONBody encodes the request body with JSON encoder.
func encodeBodyToJSON(body any) (io.Reader, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}
