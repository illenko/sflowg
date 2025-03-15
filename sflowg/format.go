package sflowg

import (
	"regexp"
	"strings"
)

var (
	hyphenStartOrEndRe = regexp.MustCompile(`(^|[^ ])-([^ ]|$)`)
	hyphenMiddleRe     = regexp.MustCompile(`([^ ])-([^ ])`)
)

func FormatKey(key string) string {
	key = strings.ReplaceAll(key, ".", "_")
	key = hyphenStartOrEndRe.ReplaceAllString(key, "${1}_${2}")
	key = hyphenMiddleRe.ReplaceAllString(key, "${1}_${2}")
	return key
}

func FormatExpression(e string) string {
	result := []rune(e)
	openParentheses := 0

	for i, r := range result {
		switch r {
		case '(':
			openParentheses++
		case ')':
			openParentheses--
		case '.':
			if openParentheses == 0 {
				result[i] = '_'
			}
		case '-':
			if openParentheses == 0 {
				temp := string(result[i-1 : i+2])
				if i == 0 {
					temp = string(result[i : i+2])
				} else if i == len(result)-1 {
					temp = string(result[i-1 : i+1])
				}

				if hyphenStartOrEndRe.MatchString(temp) || hyphenMiddleRe.MatchString(temp) {
					result[i] = '_'
				}

			}
		}
	}
	return string(result)
}
