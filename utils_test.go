package request

import "testing"

func TestGetContentType(t *testing.T) {
	if ct := getContentType(""); ct != "" {
		t.Errorf("Expect content type is \"\", actually \"%s\"", ct)
	}

	if ct := getContentType("application/json"); ct != "application/json" {
		t.Errorf("Expect content type is \"application/json\", actually \"%s\"", ct)
	}

	if ct := getContentType("application/json; charset=utf8"); ct != "application/json" {
		t.Errorf("Expect content type is \"application/json\", actually \"%s\"", ct)
	}
}
