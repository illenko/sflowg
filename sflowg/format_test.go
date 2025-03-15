package sflowg

import (
	"testing"
)

func TestFormatKey(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"Content-Type", "Content_Type"},
		{"X-Request-ID", "X_Request_ID"},
		{"User.Agent", "User_Agent"},
		{"Accept-Language", "Accept_Language"},
		{"my-custom-header", "my_custom_header"},
		{"a-b-c-d-e", "a_b_c_d_e"},
		{"request.headers.X-API-Key", "request_headers_X_API_Key"},
	}

	for _, tc := range testCases {
		actual := FormatKey(tc.input)
		if actual != tc.expected {
			t.Errorf("FormatKey(%q) = %q, expected %q", tc.input, actual, tc.expected)
		}
	}
}

func TestFormatExpression(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"a.b", "a_b"},
		{"A.B", "A_B"},
		{"a.B.c.D.E", "a_B_c_D_E"},
		{"none(tweets, {.Size > 280})", "none(tweets, {.Size > 280})"},
		{"map(tweets, {.Size})", "map(tweets, {.Size})"},
		{"request.headers.X-API-Key", "request_headers_X_API_Key"},
	}

	for _, tc := range testCases {
		result := FormatExpression(tc.input)
		if result != tc.expected {
			t.Errorf("FormatExpression(%q) = %q; want %q", tc.input, result, tc.expected)
		}
	}
}
