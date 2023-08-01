package request

import "errors"

var (
	// ErrNoURL throws when no uri and base url set in the request.
	ErrNoURL error = errors.New("no url")

	// ErrInvalidResp throws when no valid response for wrapper function.
	ErrInvalidResp error = errors.New("no response")
)
