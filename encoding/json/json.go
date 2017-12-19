package json

import (
	"encoding/json"
	cook_util "gitlab.niceprivate.com/golang/cook/util"
	"io"
	"reflect"
	"strings"
)

func Unmarshal_into_obj(j string, v interface{}) error {
	r_v := reflect.ValueOf(v)
	if r_v.Kind() != reflect.Ptr || r_v.Elem().Kind() != reflect.Struct {
		return &json.InvalidUnmarshalError{Type: r_v.Type()}
	}
	r_v = r_v.Elem()

	mapping := make(map[string]reflect.Value)
	rt_v := r_v.Type()
	for i := 0; i < rt_v.NumField(); i++ {
		rt_f := rt_v.Field(i)
		tag, ok := rt_f.Tag.Lookup("json")
		if !ok || len(tag) < 1 {
			mapping[cook_util.Hump_to_underline(rt_f.Name)] = r_v.Field(i)
		} else if tag == "-" {
			continue
		} else if tag == "-," {
			mapping["-"] = r_v.Field(i)
		} else if tag == ",omitempty" {
			mapping[cook_util.Hump_to_underline(rt_f.Name)] = r_v.Field(i)
		} else if strings.HasSuffix(tag, ",omitempty") {
			mapping[strings.TrimRight(tag, ",omitempty")] = r_v.Field(i)
		} else {
			mapping[cook_util.Hump_to_underline(rt_f.Name)] = r_v.Field(i)
		}
	}

	dec := json.NewDecoder(strings.NewReader(j))

	if t, err := dec.Token(); err == io.EOF {
		return &json.SyntaxError{Offset: 0}
	} else if err != nil {
		return err
	} else if v, ok := t.(json.Delim); !ok || v.String() != "{" {
		return &json.SyntaxError{Offset: 0}
	}

	for dec.More() {
		if t, err := dec.Token(); err == io.EOF {
			return &json.SyntaxError{Offset: 0}
		} else if err != nil {
			return err
		} else if v, ok := t.(json.Delim); ok && v.String() != "}" {
			break
		} else {
			json_col := cook_util.As_string(t)
			json_val, err := dec.Token()
			if err != nil {
				return &json.SyntaxError{Offset: 0}
			}
			if r_f, ok := mapping[json_col]; ok {
				switch r_f.Kind() {
				case reflect.Int:
					r_f.Set(reflect.ValueOf(cook_util.As_int(json_val)))
				case reflect.Int8:
					r_f.Set(reflect.ValueOf(cook_util.As_int8(json_val)))
				case reflect.Int16:
					r_f.Set(reflect.ValueOf(cook_util.As_int16(json_val)))
				case reflect.Int32:
					r_f.Set(reflect.ValueOf(cook_util.As_int32(json_val)))
				case reflect.Int64:
					r_f.Set(reflect.ValueOf(cook_util.As_int64(json_val)))
				case reflect.Uint:
					r_f.Set(reflect.ValueOf(cook_util.As_uint(json_val)))
				case reflect.Uint8:
					r_f.Set(reflect.ValueOf(cook_util.As_uint8(json_val)))
				case reflect.Uint16:
					r_f.Set(reflect.ValueOf(cook_util.As_uint16(json_val)))
				case reflect.Uint32:
					r_f.Set(reflect.ValueOf(cook_util.As_uint32(json_val)))
				case reflect.Uint64:
					r_f.Set(reflect.ValueOf(cook_util.As_uint64(json_val)))
				case reflect.String:
					r_f.Set(reflect.ValueOf(cook_util.As_string(json_val)))
				case reflect.Float32:
					r_f.Set(reflect.ValueOf(cook_util.As_float32(json_val)))
				case reflect.Float64:
					r_f.Set(reflect.ValueOf(cook_util.As_float64(json_val)))
				case reflect.Bool:
					r_f.Set(reflect.ValueOf(cook_util.As_bool(json_val)))
				}
			}
		}
	}
	return nil
}
