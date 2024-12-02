package esutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
)

// - Block

type Block struct {
	cb cipher.Block
}

// CbcBlock key长度 aes=16,24,32, des=8 ;  iv长度 des=8 aes=16
func (b *Block) CbcBlock(iv []byte, enc bool) (cipher.Block, error) {
	if len(iv) != b.cb.BlockSize() {
		return nil, ErrIvLengthError
	}
	if enc {
		return newModeBlock(cipher.NewCBCEncrypter(b.cb, iv)), nil
	}
	return newModeBlock(cipher.NewCBCDecrypter(b.cb, iv)), nil
}

func (b *Block) Block() cipher.Block {
	return b.cb
}

func aesBlock(key []byte) (*Block, error) {
	b, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return &Block{cb: b}, nil
}

// desBlock des 8
func desBlock(key []byte) (*Block, error) {
	b, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return &Block{cb: b}, nil
}

// - modelBlock

type modelBlock struct {
	mode cipher.BlockMode
}

func newModeBlock(m cipher.BlockMode) *modelBlock {
	return &modelBlock{
		mode: m,
	}
}
func (cbc *modelBlock) BlockSize() int {
	return cbc.mode.BlockSize()
}
func (cbc *modelBlock) Encrypt(dst, src []byte) {
	cbc.mode.CryptBlocks(dst, src)
}
func (cbc *modelBlock) Decrypt(dst, src []byte) {
	cbc.mode.CryptBlocks(dst, src)
}
