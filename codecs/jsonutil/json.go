package jsonutil

import "encoding/json"

func Encode(v interface{}) (string, error) {
	if v1, ok := v.(string); ok {
		return v1, nil
	}
	b, err := json.Marshal(v)
	return string(b), err
}

func Decode(d []byte, v interface{}) error {
	return json.Unmarshal(d, &v)
}
