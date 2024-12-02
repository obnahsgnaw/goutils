package md5util

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(input string) string {
	h := md5.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}
