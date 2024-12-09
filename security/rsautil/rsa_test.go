package rsautil

import (
	"crypto"
	"github.com/obnahsgnaw/goutils/security/coder"
	"testing"
)

func TestRsa(t *testing.T) {
	rsa := New(PKCS1Public(), PKCS1Private(), Encoder(coder.B64StdEncoding), SignHash(crypto.SHA256))
	priKey, pubKey, err := rsa.Generate(2048)
	if err != nil {
		t.Error(err)
		return
	}

	println("private key", string(priKey))
	println("public key", string(pubKey))

	raw := []byte("hello world")
	for i := 0; i < 30; i++ {
		raw = append(raw, []byte("hello world")...)
	}
	println("raw data:", string(raw))
	println(len(raw))

	encrypted, err := rsa.Encrypt(raw, pubKey, true)
	if err != nil {
		t.Error(err)
		return
	}
	println("encrypted data:", string(encrypted))

	decrypted, err := rsa.Decrypt(encrypted, priKey, true)
	if err != nil {
		t.Error(err)
		return
	}
	println("decrypted data:", string(decrypted))

	sign, err := rsa.Sign(raw, priKey, true)
	if err != nil {
		t.Error(err)
		return
	}
	println("sign data:", string(sign))

	err = rsa.Verify(raw, sign, pubKey, true)
	if err != nil {
		t.Error(err)
		return
	}
	println("verify ok")
}
