package util

import (
	"reflect"
)

func Assign_copy(dst, src reflect.Value) {
	switch dst.Kind() {
	case reflect.Bool:
		dst.Set(reflect.ValueOf(VAs_bool(src)))
	case reflect.Int:
		dst.Set(reflect.ValueOf(VAs_int(src)))
	case reflect.Int8:
		dst.Set(reflect.ValueOf(VAs_int8(src)))
	case reflect.Int16:
		dst.Set(reflect.ValueOf(VAs_int16(src)))
	case reflect.Int32:
		dst.Set(reflect.ValueOf(VAs_int32(src)))
	case reflect.Int64:
		dst.Set(reflect.ValueOf(VAs_int64(src)))
	case reflect.Uint:
		dst.Set(reflect.ValueOf(VAs_uint(src)))
	case reflect.Uint8:
		dst.Set(reflect.ValueOf(VAs_uint8(src)))
	case reflect.Uint16:
		dst.Set(reflect.ValueOf(VAs_uint16(src)))
	case reflect.Uint32:
		dst.Set(reflect.ValueOf(VAs_uint32(src)))
	case reflect.Uint64:
		dst.Set(reflect.ValueOf(VAs_uint64(src)))
	case reflect.Float32:
		dst.Set(reflect.ValueOf(VAs_float32(src)))
	case reflect.Float64:
		dst.Set(reflect.ValueOf(VAs_float64(src)))
	case reflect.String:
		dst.Set(reflect.ValueOf(VAs_string(src)))
	case reflect.Array:
	case reflect.Slice:
	}
}

func Is_byte_slice(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Slice && reflect.TypeOf(v).Elem().Kind() == reflect.Uint8
}
