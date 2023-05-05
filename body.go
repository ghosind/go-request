package request

import (
	"bytes"
	"encoding/json"
	"io"
)

func (cli *Client) getRequestBody(opt RequestOptions) (io.Reader, error) {
	body := opt.Body

	if body == nil {
		return nil, nil
	}

	contentType := getContentType(opt.ContentType)

	switch contentType {
	case "application/json", "": // Default JSON
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		return bytes.NewReader(data), nil
	default:
		return nil, ErrUnsupportedContentType
	}
}
