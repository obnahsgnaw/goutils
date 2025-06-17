package jsonutil

import (
	"google.golang.org/protobuf/encoding/protowire"
	"log"
	"testing"
)

func TestField(t *testing.T) {
	b, _ := Encode(map[string]interface{}{"a": 1, "b": 1})
	log.Println(IsFieldPresent([]byte(b), "a"))
	log.Println(IsFieldPresent([]byte(b), "c"))
	log.Println(IsFieldsPresent([]byte(b), map[protowire.Number]string{1: "a", 3: "c"}))
}
