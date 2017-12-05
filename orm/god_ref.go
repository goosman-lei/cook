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
	if r_model := reflect.TypeOf(m); r_model.Kind() == reflect.Ptr && r_model.Elem().Kind() == reflect.Struct && r_model.Elem().Implements(reflect.TypeOf((*Model)(nil)).Elem()) {
		model.R_Type = r_model.Elem()
	} else {
		cook_util.Panicf("model must be ptr of struct. and must implement orm.Model")
	}

	model.Mapping_with_index = make(map[int]*Ref_Field)
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
