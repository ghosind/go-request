package request

import "errors"

var (
	// ErrInvalidMethod throws when the method of the request is not a valid value.
	ErrInvalidMethod error = errors.New("invalid HTTP method")

	// ErrInvalidResp throws when no valid response for wrapper function.
	ErrInvalidResp error = errors.New("invalid response")

	// ErrNoURL throws when no uri and base url set in the request.
	ErrNoURL error = errors.New("no url")
)
