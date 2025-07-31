package structutil

import "reflect"

type NamedStruct struct {
	name string
}

func (n *NamedStruct) GetName() string {
	return n.name
}

func (n *NamedStruct) ParseName(impl interface{}) {
	tp := reflect.TypeOf(impl)
	if tp.Kind() == reflect.Ptr {
		n.name = tp.Elem().Name()
	} else {
		n.name = tp.Name()
	}
}
