package util

import "strings"

var escapes = []string{"&", "#"}

// EscapeText escapes various text sequences to make them LaTeX valid
func EscapeText(s string) string {
	for _, e := range escapes {
		s = strings.ReplaceAll(s, e, "\\"+e)
	}
	return s
}
