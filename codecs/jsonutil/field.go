package jsonutil

import "google.golang.org/protobuf/encoding/protowire"

func IsFieldPresent(data []byte, field string) bool {
	var obj map[string]interface{}
	if err := Decode(data, &obj); err == nil {
		_, ok := obj[field]
		return ok
	}
	return false
}

func IsFieldsPresent(data []byte, fields map[protowire.Number]string) map[protowire.Number]bool {
	var result = make(map[protowire.Number]bool)
	for k := range fields {
		result[k] = false
	}
	var obj map[string]interface{}
	if err := Decode(data, &obj); err == nil {
		for k, field := range fields {
			_, ok := obj[field]
			result[k] = ok
		}
	}
	return result
}
