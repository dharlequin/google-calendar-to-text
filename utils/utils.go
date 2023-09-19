package utils

import "strings"

func ExtractStringValue(line string, target string) string {
	return strings.ReplaceAll(line, target, "")
}