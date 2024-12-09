package hsutil

import (
	"crypto"
	_ "crypto/md5"
	_ "crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha512"
	"encoding/hex"
)

// Hash cal hash
func Hash(data []byte, hs crypto.Hash) (hashed []byte, err error) {
	hash := hs.New()

	if _, err = hash.Write(data); err != nil {
		return
	}

	return []byte(hex.EncodeToString(hash.Sum(nil))), nil
}

func Md5(data []byte) (hashed []byte, err error) {
	return Hash(data, crypto.MD5)
}

func Sha1(data []byte) (hashed []byte, err error) {
	return Hash(data, crypto.SHA1)
}

func Sha256(data []byte) (hashed []byte, err error) {
	return Hash(data, crypto.SHA256)
}

func Sha512(data []byte) (hashed []byte, err error) {
	return Hash(data, crypto.SHA512)
}
