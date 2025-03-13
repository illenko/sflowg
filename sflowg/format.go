package sflowg

import (
	"regexp"
	"strings"
)

func FormatKey(k string) string {
	return strings.ReplaceAll(k, ".", "_")
}

var re = regexp.MustCompile(`\.`)

func FormatExpression(e string) string {
	matches := re.FindAllStringIndex(e, -1)

	result := []rune(e)
	for _, match := range matches {
		dotIndex := match[0]
		openParentheses := 0

		for i := 0; i < dotIndex; i++ {
			if result[i] == '(' {
				openParentheses++
			} else if result[i] == ')' {
				openParentheses--
			}
		}

		if openParentheses == 0 {
			result[dotIndex] = '_'
		}
	}
	return string(result)
}
