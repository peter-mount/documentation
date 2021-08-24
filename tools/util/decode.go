package util

import (
	"strconv"
)

func DecodeString(v interface{}, def string) string {
	var r string

	if v != nil {
		if s, ok := v.(string); ok {
			r = s
		} else if i, ok := v.(int); ok {
			r = strconv.Itoa(i)
		}
	}

	if r == "" {
		r = def
	}

	return r
}

func DecodeInt(v interface{}, def int) (int, bool) {
	r := def

	if v != nil {
		if s, ok := v.(string); ok {
			i, err := strconv.Atoi(s)
			if err != nil {
				return 0, false
			}
			r = i
		} else if i, ok := v.(int); ok {
			r = i
		}
	}

	return r, true
}
