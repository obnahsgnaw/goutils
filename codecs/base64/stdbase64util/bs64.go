package base64

import "encoding/base64"

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
