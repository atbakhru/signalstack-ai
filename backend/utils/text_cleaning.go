package utils

import (
	"regexp"
	"strings"
)

var whitespaceNormalizer = regexp.MustCompile(`\s+`)

func CleanText(text string) string {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return ""
	}
	return whitespaceNormalizer.ReplaceAllString(trimmed, " ")
}
