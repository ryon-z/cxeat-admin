package common

import (
	"reflect"
)

// ChkPtr nil 체크
func ChkPtr(given interface{}) interface{} {
	if reflect.TypeOf(given).Kind() != reflect.Ptr {
		return given
	}

	if !reflect.ValueOf(given).IsNil() {
		return reflect.ValueOf(given).Elem().Interface()
	}

	switch reflect.TypeOf(given).Elem().Kind() {
	default:
		return given
	case reflect.String:
		return ""
	case reflect.Bool:
		return false
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return 0
	case reflect.Float32, reflect.Float64:
		return 0
	}
}
