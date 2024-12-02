package esutil

import (
	"crypto/aes"
	"crypto/des"
	"github.com/obnahsgnaw/goutils/randutil"
)

const (
	Des    EsType = 8
	Aes128 EsType = 16
	Aes192 EsType = 24
	Aes256 EsType = 32
)

//AES-128：key长度16 字节
//AES-192：key长度24 字节
//AES-256：key长度32 字节

type EsType int

func (e EsType) KeyLen() int {
	return int(e)
}

func (e EsType) IvLen() int {
	if e == Des {
		return des.BlockSize
	}

	return aes.BlockSize
}

func (e EsType) RandIv() []byte {
	return []byte(randutil.RandNum(e.IvLen()))
}

func (e EsType) RandKey() []byte {
	return []byte(randutil.RandAlpha(e.KeyLen()))
}
