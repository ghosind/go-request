package request

import "strings"

// getContentType gets the first part of the value of the `Content-Type` field in the response
// headers.
func getContentType(s string) string {
	contentType, _, _ := strings.Cut(s, ";") // remove parameters

	return strings.TrimSpace(contentType)
}
