package protoutil

import "google.golang.org/protobuf/encoding/protowire"

func IsFieldPresent(data []byte, fieldNumber protowire.Number) bool {
	for len(data) > 0 {
		tag, wireType, n := protowire.ConsumeTag(data)
		if n < 0 {
			return false
		}
		if tag == fieldNumber {
			return true
		}

		m := protowire.ConsumeFieldValue(tag, wireType, data[n:])
		if m < 0 {
			return false
		}
		data = data[n+m:]
	}
	return false
}

func IsFieldsPresent(data []byte, fields []protowire.Number) (result map[protowire.Number]bool) {
	result = make(map[protowire.Number]bool)
	for _, field := range fields {
		result[field] = false
	}
	for len(data) > 0 {
		tag, wireType, n := protowire.ConsumeTag(data)
		if n < 0 {
			return
		}
		if _, ok := result[tag]; ok {
			result[tag] = true
		}
		m := protowire.ConsumeFieldValue(tag, wireType, data[n:])
		if m < 0 {
			return
		}
		data = data[n+m:]
	}
	return
}
