package request

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

const (
	// RequestContentTypeJSON indicates the request body encodes as a JSON.
	RequestContentTypeJSON string = "json"
)

// getRequestBody returns the encoded request body as an io.Reader object. The function will try
// to get a supported content type from the Header field in the request config, and it will try to
// serialize the data as a JSON string if no content type or the content type is unsupported.
// It'll skip encoding the request body if it's a nil pointer, a string, or a slice of bytes.
func (cli *Client) getRequestBody(opt RequestOptions) (io.Reader, error) {
	body := opt.Body
	if body == nil {
		return nil, nil
	}

	data, err := cli.encodeRequestBody(body, opt.ContentType)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}

// encodeRequestBody encodes the request body by the specific encoder. It'll get the encoding
// type from the 'ContentType' field in the request config, and it uses JSON as the default
// encoder. It will skip encoding the request body data if it is a byte array slice or a string.
func (cli *Client) encodeRequestBody(body any, contentType string) ([]byte, error) {
	if body == nil {
		return nil, nil
	}

	// Encoding byte slice or string is unnecessary.
	if v, ok := body.([]byte); ok {
		return v, nil
	} else if v, ok := body.(string); ok {
		return []byte(v), nil
	}

	var handler func(any) ([]byte, error)

	switch strings.ToLower(contentType) {
	case RequestContentTypeJSON, "":
		handler = json.Marshal
	default:
		return nil, ErrUnsupportedType
	}

	return handler(body)
}
