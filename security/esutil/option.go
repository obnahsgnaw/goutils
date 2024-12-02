package esutil

import "github.com/obnahsgnaw/goutils/security/coder"

type Option func(*ADes)

func Encoder(s coder.Encoder) Option {
	return func(e *ADes) {
		if s != nil {
			e.encoder = s
		}
	}
}
