package strutil

import "strings"

func PadLen(str string, max int) string {
	return PadLenWith(str, max, " ")
}

func PadLenWith(str string, max int, rp string) string {
	if sp := max - len(str); sp > 0 {
		str = str + strings.Repeat(rp, sp)
	}
	return str
}
