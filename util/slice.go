package util

import (
	"reflect"
	"sort"
	"strconv"
)

func Slice_uniq_string(s sort.StringSlice) []string {
	t := []string{}
	if s.Len() < 1 {
		return t
	}
	s.Sort()
	o := 0
	t = append(t, s[o])
	for i := 1; i < s.Len(); i++ {
		if s[i] != t[o] {
			t = append(t, s[i])
			o++
		}
	}
	return t
}

func Slice_seqs(start, stop, step int) []int {
	s := []int{}
	for i := start; i <= stop; i += step {
		s = append(s, i)
	}
	return s
}

func Slice_int_to_string(s []int) []string {
	sr := make([]string, len(s))
	for i, v := range s {
		sr[i] = strconv.FormatInt(int64(v), 10)
	}
	return sr
}

func Slice_pick(s interface{}, field string, fn func(v reflect.Value)) {
	r_type_s := reflect.TypeOf(s)
	r_value_s := reflect.ValueOf(s)
	if r_type_s.Kind() == reflect.Ptr {
		r_type_s = r_type_s.Elem()
		r_value_s = r_value_s.Elem()
	}
	if r_type_s.Kind() != reflect.Slice {
		return
	}

	for i := 0; i < r_value_s.Len(); i++ {
		r_value_e := r_value_s.Index(i)
		if r_value_e.Kind() == reflect.Interface {
			r_value_e = r_value_e.Elem()
		}
		if r_value_e.Kind() == reflect.Ptr {
			r_value_e = r_value_e.Elem()
		}
		fn(r_value_e.FieldByName(field))
	}
}

func Slice_pick_int(s interface{}, field string) []int {
	r_ints := []int{}
	Slice_pick(s, field, func(v reflect.Value) {
		r_ints = append(r_ints, VAs_int(v))
	})
	return r_ints
}

func Slice_pick_int8(s interface{}, field string) []int8 {
	r_int8s := []int8{}
	Slice_pick(s, field, func(v reflect.Value) {
		r_int8s = append(r_int8s, VAs_int8(v))
	})
	return r_int8s
}

func Slice_pick_int16(s interface{}, field string) []int16 {
	r_int16s := []int16{}
	Slice_pick(s, field, func(v reflect.Value) {
		r_int16s = append(r_int16s, VAs_int16(v))
	})
	return r_int16s
}

func Slice_pick_int32(s interface{}, field string) []int32 {
	r_int32s := []int32{}
	Slice_pick(s, field, func(v reflect.Value) {
		r_int32s = append(r_int32s, VAs_int32(v))
	})
	return r_int32s
}

func Slice_pick_int64(s interface{}, field string) []int64 {
	r_int64s := []int64{}
	Slice_pick(s, field, func(v reflect.Value) {
		r_int64s = append(r_int64s, VAs_int64(v))
	})
	return r_int64s
}

func Slice_pick_uint(s interface{}, field string) []uint {
	r_uints := []uint{}
	Slice_pick(s, field, func(v reflect.Value) {
		r_uints = append(r_uints, VAs_uint(v))
	})
	return r_uints
}

func Slice_pick_uint8(s interface{}, field string) []uint8 {
	r_uint8s := []uint8{}
	Slice_pick(s, field, func(v reflect.Value) {
		r_uint8s = append(r_uint8s, VAs_uint8(v))
	})
	return r_uint8s
}

func Slice_pick_uint16(s interface{}, field string) []uint16 {
	r_uint16s := []uint16{}
	Slice_pick(s, field, func(v reflect.Value) {
		r_uint16s = append(r_uint16s, VAs_uint16(v))
	})
	return r_uint16s
}

func Slice_pick_uint32(s interface{}, field string) []uint32 {
	r_uint32s := []uint32{}
	Slice_pick(s, field, func(v reflect.Value) {
		r_uint32s = append(r_uint32s, VAs_uint32(v))
	})
	return r_uint32s
}

func Slice_pick_uint64(s interface{}, field string) []uint64 {
	r_uint64s := []uint64{}
	Slice_pick(s, field, func(v reflect.Value) {
		r_uint64s = append(r_uint64s, VAs_uint64(v))
	})
	return r_uint64s
}

func Slice_pick_bool(s interface{}, field string) []bool {
	r_bools := []bool{}
	Slice_pick(s, field, func(v reflect.Value) {
		r_bools = append(r_bools, VAs_bool(v))
	})
	return r_bools
}

func Slice_pick_float32(s interface{}, field string) []float32 {
	r_float32s := []float32{}
	Slice_pick(s, field, func(v reflect.Value) {
		r_float32s = append(r_float32s, VAs_float32(v))
	})
	return r_float32s
}

func Slice_pick_float64(s interface{}, field string) []float64 {
	r_float64s := []float64{}
	Slice_pick(s, field, func(v reflect.Value) {
		r_float64s = append(r_float64s, VAs_float64(v))
	})
	return r_float64s
}

func Slice_pick_string(s interface{}, field string) []string {
	r_strings := []string{}
	Slice_pick(s, field, func(v reflect.Value) {
		r_strings = append(r_strings, VAs_string(v))
	})
	return r_strings
}

func Slice_string_fill_bool(s []string, v bool) map[string]bool {
	m := make(map[string]bool, len(s))
	for _, t := range s {
		m[t] = v
	}
	return m
}

func Slice_string_remove(subject []string, targets ...string) []string {
	if len(targets) < 1 {
		return subject
	}
	m := Slice_string_fill_bool(targets, true)
	return Slice_string_filter(subject, func(s string) bool {
		_, ok := m[s]
		return ok
	})
}

func Slice_string_filter(subject []string, fn func(string) bool, n ...uint) []string {
	var c uint = 0
	var limit uint = 0
	if len(n) > 0 {
		limit = n[0]
	}
	var newSubject = []string{}
	for i, s := range subject {
		if !fn(s) {
			newSubject = append(newSubject, s)
		} else {
			c++
			if c == limit {
				newSubject = append(newSubject, subject[i+1:]...)
				break
			}
		}
	}
	return newSubject
}

func Slice_string_remove_n(subject []string, n uint, targets ...string) []string {
	if len(targets) < 1 || n == 0 {
		return subject
	}
	m := Slice_string_fill_bool(targets, true)
	return Slice_string_filter(subject, func(s string) bool {
		_, ok := m[s]
		return ok
	}, n)
}
