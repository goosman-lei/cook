package util

import (
	"fmt"
	"reflect"
	"strconv"
)

func Indirect(v interface{}) reflect.Value {
	return V_indirect(reflect.ValueOf(v))
}

func Indirect_zero(v interface{}, wantZero bool) reflect.Value {
	return V_indirect_zero(reflect.ValueOf(v), wantZero)
}

func Indirect_stopAt(v interface{}, stopFn func(reflect.Value) bool) reflect.Value {
	return V_indirect_stopAt(reflect.ValueOf(v), stopFn)
}

// Copy from GOROOT/src/encoding/json/decode.go: func (d *decodeState) indirect()
func Indirect_full(v interface{}, stopFn reflect.Value, wantZero bool) reflect.Value {
	return V_indirect_full(reflect.ValueOf(v), stopFn, wantZero)
}

func V_indirect(v reflect.Value) reflect.Value {
	return V_indirect_full(v, reflect.ValueOf(func(reflect.Value) bool {
		return false
	}), false)
}

func V_indirect_zero(v reflect.Value, wantZero bool) reflect.Value {
	return V_indirect_full(v, reflect.ValueOf(func(reflect.Value) bool {
		return false
	}), wantZero)
}

func V_indirect_stopAt(v reflect.Value, stopFn func(reflect.Value) bool) reflect.Value {
	return V_indirect_full(v, reflect.ValueOf(stopFn), false)
}

// Copy from GOROOT/src/encoding/json/decode.go: func (d *decodeState) indirect()
func V_indirect_full(v reflect.Value, stopFn reflect.Value, wantZero bool) (rv reflect.Value) {
	defer func() {
		if e := recover(); e != nil {
			rv = reflect.ValueOf(nil)
		}
	}()
	// make v CanAddr()
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}
	for {
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && (!wantZero || e.Elem().Kind() == reflect.Ptr) {
				v = e
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if v.Elem().Kind() != reflect.Ptr && wantZero && v.CanSet() {
			break
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		// stopFn decide This is valid Type
		if stopFn.Kind() == reflect.Func {
			out := stopFn.Call([]reflect.Value{reflect.ValueOf(v)})
			if len(out) > 0 && VAs_bool(out[0]) {
				return v
			}
		}
		v = v.Elem()
	}
	return v
}

func Is_empty(v interface{}) bool {
	return VIs_empty(reflect.ValueOf(v))
}

func VIs_empty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func PAs_string(v interface{}) string {
	return VAs_string(reflect.Indirect(reflect.ValueOf(v)))
}

func PAs_int(v interface{}) int {
	return VAs_int(reflect.Indirect(reflect.ValueOf(v)))
}

func PAs_int8(v interface{}) int8 {
	return VAs_int8(reflect.Indirect(reflect.ValueOf(v)))
}

func PAs_int16(v interface{}) int16 {
	return VAs_int16(reflect.Indirect(reflect.ValueOf(v)))
}

func PAs_int32(v interface{}) int32 {
	return VAs_int32(reflect.Indirect(reflect.ValueOf(v)))
}

func PAs_int64(v interface{}) int64 {
	return VAs_int64(reflect.Indirect(reflect.ValueOf(v)))
}

func PAs_uint(v interface{}) uint {
	return VAs_uint(reflect.Indirect(reflect.ValueOf(v)))
}

func PAs_uint8(v interface{}) uint8 {
	return VAs_uint8(reflect.Indirect(reflect.ValueOf(v)))
}

func PAs_uint16(v interface{}) uint16 {
	return VAs_uint16(reflect.Indirect(reflect.ValueOf(v)))
}

func PAs_uint32(v interface{}) uint32 {
	return VAs_uint32(reflect.Indirect(reflect.ValueOf(v)))
}

func PAs_uint64(v interface{}) uint64 {
	return VAs_uint64(reflect.Indirect(reflect.ValueOf(v)))
}

func PAs_float32(v interface{}) float32 {
	return VAs_float32(reflect.Indirect(reflect.ValueOf(v)))
}

func PAs_float64(v interface{}) float64 {
	return VAs_float64(reflect.Indirect(reflect.ValueOf(v)))
}

func PAs_bool(v interface{}) bool {
	return VAs_bool(reflect.Indirect(reflect.ValueOf(v)))
}

func As_string(v interface{}) string {
	return VAs_string(reflect.ValueOf(v))
}

func As_int(v interface{}) int {
	return VAs_int(reflect.ValueOf(v))
}

func As_int8(v interface{}) int8 {
	return VAs_int8(reflect.ValueOf(v))
}

func As_int16(v interface{}) int16 {
	return VAs_int16(reflect.ValueOf(v))
}

func As_int32(v interface{}) int32 {
	return VAs_int32(reflect.ValueOf(v))
}

func As_int64(v interface{}) int64 {
	return VAs_int64(reflect.ValueOf(v))
}

func As_uint(v interface{}) uint {
	return VAs_uint(reflect.ValueOf(v))
}

func As_uint8(v interface{}) uint8 {
	return VAs_uint8(reflect.ValueOf(v))
}

func As_uint16(v interface{}) uint16 {
	return VAs_uint16(reflect.ValueOf(v))
}

func As_uint32(v interface{}) uint32 {
	return VAs_uint32(reflect.ValueOf(v))
}

func As_uint64(v interface{}) uint64 {
	return VAs_uint64(reflect.ValueOf(v))
}

func As_float32(v interface{}) float32 {
	return VAs_float32(reflect.ValueOf(v))
}

func As_float64(v interface{}) float64 {
	return VAs_float64(reflect.ValueOf(v))
}

func As_bool(v interface{}) bool {
	return VAs_bool(reflect.ValueOf(v))
}

func VAs_string(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'f', -1, 32)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return string(v.Bytes())
		} else {
			return fmt.Sprintf("%v", v)
		}
	default:
		return fmt.Sprintf("%v", v)
	}
}

func VAs_int(v reflect.Value) int {
	return int(VAs_int64(v))
}

func VAs_int8(v reflect.Value) int8 {
	return int8(VAs_int64(v))
}

func VAs_int16(v reflect.Value) int16 {
	return int16(VAs_int64(v))
}

func VAs_int32(v reflect.Value) int32 {
	return int32(VAs_int64(v))
}

func VAs_int64(v reflect.Value) int64 {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(v.Uint())
	case reflect.Float64, reflect.Float32:
		return int64(v.Float())
	case reflect.Bool:
		if v.Bool() {
			return 1
		} else {
			return 0
		}
	case reflect.String:
		if iv, err := strconv.ParseInt(v.String(), 10, 64); err != nil {
			return 0
		} else {
			return iv
		}
	default:
		return 0
	}
}

func VAs_uint(v reflect.Value) uint {
	return uint(VAs_uint64(v))
}

func VAs_uint8(v reflect.Value) uint8 {
	return uint8(VAs_uint64(v))
}

func VAs_uint16(v reflect.Value) uint16 {
	return uint16(VAs_uint64(v))
}

func VAs_uint32(v reflect.Value) uint32 {
	return uint32(VAs_uint64(v))
}

func VAs_uint64(v reflect.Value) uint64 {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uint64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uint64(v.Uint())
	case reflect.Float64, reflect.Float32:
		return uint64(v.Float())
	case reflect.Bool:
		if v.Bool() {
			return 1
		} else {
			return 0
		}
	case reflect.String:
		if iv, err := strconv.ParseUint(v.String(), 10, 64); err != nil {
			return 0
		} else {
			return iv
		}
	default:
		return 0
	}
}

func VAs_float32(v reflect.Value) float32 {
	return float32(VAs_float64(v))
}

func VAs_float64(v reflect.Value) float64 {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint())
	case reflect.Float64, reflect.Float32:
		return v.Float()
	case reflect.Bool:
		if v.Bool() {
			return 1.0
		} else {
			return 0.0
		}
	case reflect.String:
		if fv, err := strconv.ParseFloat(v.String(), 64); err != nil {
			return 0.0
		} else {
			return fv
		}
	default:
		return 0.0
	}
}

func VAs_bool(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() != 0
	case reflect.Float64, reflect.Float32:
		return v.Float() != 0.0
	case reflect.Bool:
		return v.Bool()
	case reflect.String, reflect.Array:
		return v.Len() > 0
	case reflect.Ptr, reflect.Func, reflect.Interface:
		return !v.IsNil()
	case reflect.Chan, reflect.Map, reflect.Slice:
		return !v.IsNil() && v.Len() > 0
	default:
		return false
	}
}
