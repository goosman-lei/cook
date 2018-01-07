package json

import (
	"bytes"
	"encoding/json"
	cook_util "gitlab.niceprivate.com/golang/cook/util"
	"io"
	"reflect"
	"strings"
)

func Marshal_obj(v interface{}) (string, error) {
	r_v := reflect.ValueOf(v)
	if r_v.Kind() != reflect.Ptr || r_v.Elem().Kind() != reflect.Struct {
		return "", &json.InvalidUnmarshalError{Type: r_v.Type()}
	}
	r_v = r_v.Elem()

	buf := new(bytes.Buffer)
	buf.WriteRune('{')
	is_first := true
	rt_v := r_v.Type()
	for i := 0; i < rt_v.NumField(); i++ {
		rt_f := rt_v.Field(i)
		tag, ok := rt_f.Tag.Lookup("json")
		if tag == "-" || (strings.HasSuffix(tag, ",omitempty") && cook_util.VIs_empty(r_v.Field(i))) {
			continue
		}
		if is_first {
			is_first = false
		} else {
			buf.WriteRune(',')
		}
		if !ok || len(tag) < 1 || tag == ",omitempty" {
			buf.WriteRune('"')
			buf.WriteString(cook_util.Hump_to_underline(rt_f.Name))
			buf.WriteRune('"')
			buf.WriteRune(':')
		} else if tag == "-," {
			buf.WriteString("\"-\":")
		} else if strings.HasSuffix(tag, ",omitempty") {
			buf.WriteRune('"')
			buf.WriteString(strings.TrimRight(tag, ",omitempty"))
			buf.WriteRune('"')
			buf.WriteRune(':')
		} else {
			buf.WriteRune('"')
			buf.WriteString(tag)
			buf.WriteRune('"')
			buf.WriteRune(':')
		}
		switch r_v.Field(i).Kind() {
		case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64:
			// In json, there have no 64 bit integer support. It assume that as double when encounter integer over 2^32
			// so here we convert it to string
			if j, err := json.Marshal(cook_util.VAs_string(r_v.Field(i))); err == nil {
				buf.Write(j)
			} else {
				buf.WriteString("null")
			}
		default:
			if j, err := json.Marshal(r_v.Field(i).Interface()); err == nil {
				buf.Write(j)
			} else {
				buf.WriteString("null")
			}
		}
	}
	buf.WriteRune('}')
	return buf.String(), nil
}
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
			mapping[tag] = r_v.Field(i)
		}
	}

	dec := json.NewDecoder(strings.NewReader(j))
	dec.UseNumber()

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
