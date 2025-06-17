package jsonutil

func IsFieldPresent(data []byte, field string) bool {
	var obj map[string]interface{}
	if err := Decode(data, &obj); err == nil {
		_, ok := obj[field]
		return ok
	}
	return false
}

func IsFieldsPresent(data []byte, fields []string) map[string]bool {
	var result = make(map[string]bool)
	for _, field := range fields {
		result[field] = false
	}
	var obj map[string]interface{}
	if err := Decode(data, &obj); err == nil {
		for _, field := range fields {
			_, ok := obj[field]
			result[field] = ok
		}
	}
	return result
}
