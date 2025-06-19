package structutil

import (
	"fmt"
	"google.golang.org/protobuf/runtime/protoimpl"
	"testing"
)

func TestStruct2map(t *testing.T) {
	d := &Config{
		Debug: false,
	}

	m := Struct2map(d)

	fmt.Printf("%v \n", m)
}

type Config struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Debug bool `protobuf:"varint,1,opt,name=debug,proto3" json:"debug,omitempty"`
}
