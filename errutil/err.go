package errutil

import (
	"errors"
	"fmt"
	"github.com/obnahsgnaw/goutils/strutil"
)

func New(s ...string) error {
	return errors.New(strutil.ToString(s...))
}

func NewFromError(err error, s ...string) error {
	if err == nil {
		return New(s...)
	}

	return fmt.Errorf(strutil.ToString(s...)+": %w", err)
}
