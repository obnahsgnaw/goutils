package esutil

import (
	"github.com/obnahsgnaw/goutils/security/coder"
	"testing"
)

func TestEs(t *testing.T) {
	es := New(Aes256, CbcMode, Encoder(coder.B64StdEncoding))

	raw := []byte("hello world")
	for i := 0; i < 30; i++ {
		raw = append(raw, []byte("hello world")...)
	}
	println("raw:", string(raw))
	println(len(raw))

	key := Aes256.RandKey()

	encrypted, iv, err := es.Encrypt(raw, key, true)
	if err != nil {
		t.Fatal(err)
		return
	}
	println("encrypted:", string(encrypted), string(iv))

	decrypted, err := es.Decrypt(encrypted, key, iv, true)
	if err != nil {
		t.Fatal(err)
		return
	}
	println("decrypted:", string(decrypted))
}
