package request

import "strings"

func getContentType(s string) string {
	contentType, _, _ := strings.Cut(s, ";") // remove parameters

	return strings.TrimSpace(contentType)
}
