package orm

import (
	cook_util "gitlab.niceprivate.com/golang/cook/util"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	TAG_KEY = "orm"
)

type Ref_Tag struct {
	FuncName string
	FuncArgs []interface{}
	Handler  func(*Ref_Field, *Ref_Tag)
}

type Ref_Field struct {
	Model         *Ref_Model
	R_StructField reflect.StructField
	Column        string
	Tags          map[string]*Ref_Tag
}

var (
	pattern_func        *regexp.Regexp                        = regexp.MustCompile("^(\\w+)\\(?(.*?)\\)?$")
	tag_handler_mapping map[string]func(*Ref_Field, *Ref_Tag) = map[string]func(*Ref_Field, *Ref_Tag){
		"pk":  tag_handler_pk,
		"col": tag_handler_col,
	}
)

func NewRefField(model *Ref_Model, idx int) *Ref_Field {
	field := &Ref_Field{}
	field.Model = model
	field.R_StructField = model.R_Type.Field(idx)
	field.apply_tags()
	return field
}

func (f *Ref_Field) apply_tags() {
	f.Tags = make(map[string]*Ref_Tag)
	tag_str := f.R_StructField.Tag.Get(TAG_KEY)

	func_chunks := strings.Split(tag_str, ";")
	for _, func_chunk := range func_chunks {
		matches := pattern_func.FindStringSubmatch(func_chunk)
		if len(matches) != 3 {
			continue
		}
		func_name := matches[1]
		func_args := []interface{}{}
		for _, raw_arg := range strings.Split(matches[2], ",") {
			raw_arg = strings.Trim(raw_arg, " ")
			if len(raw_arg) == 0 {
				continue
			} else if i_val, i_err := strconv.ParseInt(raw_arg, 10, 64); i_err == nil {
				func_args = append(func_args, int(i_val))
			} else if f_val, f_err := strconv.ParseFloat(raw_arg, 64); f_err == nil {
				func_args = append(func_args, float64(f_val))
			} else {
				func_args = append(func_args, strings.Trim(raw_arg, "\""))
			}
		}

		if handler, ok := tag_handler_mapping[func_name]; ok {
			f.Tags[func_name] = &Ref_Tag{
				FuncName: func_name,
				FuncArgs: func_args,
				Handler:  handler,
			}
		}
	}

	if _, ok := f.Tags["col"]; !ok {
		f.Tags["col"] = &Ref_Tag{
			FuncName: "col",
			FuncArgs: []interface{}{},
			Handler:  tag_handler_col,
		}
	}

	for _, tag := range f.Tags {
		tag.Handler(f, tag)
	}
}

func tag_handler_pk(f *Ref_Field, t *Ref_Tag) {
	if len(t.FuncArgs) == 0 {
		f.Model.PK = append(f.Model.PK, f)
	} else if len(t.FuncArgs) == 2 {
		if idx, ok := t.FuncArgs[0].(int); ok {
			if cnt, ok := t.FuncArgs[1].(int); ok {
				if len(f.Model.PK) == 0 {
					f.Model.PK = make([]*Ref_Field, cnt)
				}
				f.Model.PK[idx] = f
			}
		}
	}
}

func tag_handler_col(f *Ref_Field, t *Ref_Tag) {
	if len(t.FuncArgs) >= 1 {
		if v, ok := t.FuncArgs[0].(string); ok {
			f.Column = v
		}
	} else {
		f.Column = cook_util.Hump_to_underline(f.R_StructField.Name)
	}
}
