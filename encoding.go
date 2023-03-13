package request

import (
	"encoding/json"
	"io"
)

func encodeJson(body any) ([]byte, error) {
	return json.Marshal(body)
}

func decodeJson(body io.Reader, out any) error {
	return json.NewDecoder(body).Decode(out)
}
