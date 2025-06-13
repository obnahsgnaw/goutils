package protoutil

import (
	"google.golang.org/protobuf/encoding/protowire"
	"log"
	"testing"
)

func TestField(t *testing.T) {
	protodB := []byte{10, 12, 8, 247, 254, 174, 194, 6, 16, 232, 177, 188, 136, 3, 21, 0, 0, 246, 66}
	log.Println(IsFieldPresent(protodB, 1))
	log.Println(IsFieldPresent(protodB, 2))
	log.Println(IsFieldPresent(protodB, 3))
	log.Println(IsFieldPresent(protodB, 4))
	log.Println(IsFieldPresent(protodB, 5))
	log.Println(IsFieldPresent(protodB, 6))
	log.Println(IsFieldPresent(protodB, 7))
	log.Println(IsFieldPresent(protodB, 8))
}

func TestFields(t *testing.T) {
	protodB := []byte{10, 12, 8, 247, 254, 174, 194, 6, 16, 232, 177, 188, 136, 3, 21, 0, 0, 246, 66}
	log.Println(IsFieldsPresent(protodB, []protowire.Number{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}))
	protodc := []byte{10, 11, 8, 160, 132, 175, 194, 6, 16, 240, 253, 251, 103, 21, 0, 0, 246, 66, 53, 0, 0, 48, 65}
	log.Println(IsFieldsPresent(protodc, []protowire.Number{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}))
}
