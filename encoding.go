package request

import (
	"encoding/json"
	"io"
)

func decodeJson(body io.Reader, out any) error {
	return json.NewDecoder(body).Decode(out)
}
