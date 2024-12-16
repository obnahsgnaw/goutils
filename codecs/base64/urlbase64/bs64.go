package base64

import "encoding/base64"

func Encode(b []byte) string {
	return base64.URLEncoding.EncodeToString(b)
}

func Decode(s string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(s)
}
