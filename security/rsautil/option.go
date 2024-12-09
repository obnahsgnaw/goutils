package rsautil

import (
	"crypto"
	"github.com/obnahsgnaw/goutils/security/coder"
)

type Option func(*Rsa)

func PKCS1Private() Option {
	return func(e *Rsa) {
		e.prkType = prkPKCS1
	}
}

func PKCS8Private() Option {
	return func(e *Rsa) {
		e.prkType = prkPKCS8
	}
}

func PKCS1Public() Option {
	return func(e *Rsa) {
		e.pukType = pukPKCS1
	}
}

func PKIXPublic() Option {
	return func(e *Rsa) {
		e.pukType = pukPKIX
	}
}

func Encoder(s coder.Encoder) Option {
	return func(e *Rsa) {
		if s != nil {
			e.encoder = s
		}
	}
}

func SignHash(hs crypto.Hash) Option {
	return func(e *Rsa) {
		e.signHash = hs
	}
}
