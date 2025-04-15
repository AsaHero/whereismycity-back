package utility

import (
	"bytes"
	"encoding/json"
	"reflect"
)

// Ter - тернарный оператор
func Ter[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}

// GetTypeName returns the type Name of the given interface
func GetTypeName(i interface{}) string {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	}
	return t.Name()
}

// FormatStruct возвращает строку для отображения данных в удобно-читаемом формате
func FormatStruct(i interface{}) string {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetIndent("", "\t")
	enc.SetEscapeHTML(false)
	err := enc.Encode(i)
	if err != nil {
		return ""
	}
	return buf.String()
}
