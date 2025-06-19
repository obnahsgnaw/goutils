package structutil

import (
	"reflect"
	"strings"
	"unicode"
)

func Struct2map(s interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	v := reflect.ValueOf(s)
	if v.IsNil() {
		return nil
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !unicode.IsUpper([]rune(field.Name)[0]) {
			continue
		}
		value := v.Field(i).Interface()
		key := strings.ToLower(field.Name)
		if value == nil {
			continue
		}
		// 检查值是否也是一个结构体。
		vt := reflect.TypeOf(value)
		if vt.Kind() == reflect.Ptr {
			vt = vt.Elem()
		}
		if vt.Kind() == reflect.Struct {
			// 将嵌套结构体转换为 map。
			nestedMap := Struct2map(value)
			if nestedMap != nil {
				m[key] = nestedMap
			}
		} else {
			m[key] = value
		}
	}

	return m
}
