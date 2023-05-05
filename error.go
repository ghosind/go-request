package request

import "errors"

// ErrNoURL throws when no uri and base url set in the request.
var ErrNoURL error = errors.New("no url")

var ErrUnsupportedContentType = errors.New("unsupported content type")
