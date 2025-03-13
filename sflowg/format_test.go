package sflowg

import (
	"testing"
)

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
	}

	for _, tc := range testCases {
		result := FormatExpression(tc.input)
		if result != tc.expected {
			t.Errorf("FormatExpression(%q) = %q; want %q", tc.input, result, tc.expected)
		}
	}
}
