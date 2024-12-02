package bs64util

import "encoding/base64"

func StdEncode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func StdDecode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func UrlEncode(b []byte) string {
	return base64.URLEncoding.EncodeToString(b)
}

func UrlDecode(s string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(s)
}
