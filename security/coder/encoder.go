package coder

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
)

var (
	B64StdEncoding = base64.StdEncoding
	B64UrlEncoding = base64.URLEncoding
	HexEncoding    = HexEncoder{}
)

type Encoder interface {
	EncodeToString([]byte) string
	DecodeString(string) ([]byte, error)
}

type HexEncoder struct{}

func (HexEncoder) EncodeToString(data []byte) string {
	return strings.ToUpper(hex.EncodeToString(data))
}

func (HexEncoder) DecodeString(data string) ([]byte, error) {
	return hex.DecodeString(data)
}
