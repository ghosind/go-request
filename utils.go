package request

import "strings"

// getResponseType gets the type of the response's content from the `Content-Type` field in the
// response headers.
func getContentType(contentType string) string {
	contentType = strings.ToLower(contentType)

	switch {
	case strings.Contains(contentType, "json"):
		return RequestContentTypeJSON
	case contentType == "":
		return RequestContentTypeJSON
	default:
		return "unknown"
	}
}
