package util

import "strconv"

var units = []string{"", "K", "M", "G", "T"}

func Unit(v int) string {
	var i int
	for i = 0; v >= 1024; i++ {
		v /= 1024
	}
	return strconv.Itoa(v) + units[i]
}
