package orm

import (
	cook_util "gitlab.niceprivate.com/golang/cook/util"
	"reflect"
)

type Ref_Model struct {
	R_Type              reflect.Type
	PK                  []*Ref_Field
	Mapping_with_index  map[int]*Ref_Field
	Mapping_with_field  map[string]*Ref_Field
	Mapping_with_column map[string]*Ref_Field
}

func NewRefModel(m interface{}) *Ref_Model {
	model := &Ref_Model{}
	model.R_Type = cook_util.TypeOf_Must_struct_ptr(m)
	model.Mapping_with_index = make(map[string]*Ref_Field)
	model.Mapping_with_field = make(map[string]*Ref_Field)
	model.Mapping_with_column = make(map[string]*Ref_Field)

	for i := 0; i < model.R_Type.NumField(); i++ {
		field := NewRefField(model, i)
		if field.NoMapping {
			continue
		}
		model.Mapping_with_index[i] = field
		model.Mapping_with_field[field.R_StructField.Name] = field
		model.Mapping_with_column[field.Column] = field
	}

	return model
}
