package errutil

import (
	"errors"
	"log"
	"testing"
)

func TestErrBuilder(t *testing.T) {
	b := NewBuilder("user controller")
	err := b.New(errors.New("db disconnected"))
	if err != nil {
		log.Println(err.Error())
	}
}
