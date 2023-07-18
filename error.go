package request

import "errors"

var (
	// ErrNoURL throws when no uri and base url set in the request.
	ErrNoURL error = errors.New("no url")
)
