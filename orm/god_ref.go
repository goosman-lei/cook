package orm

import (
	cook_util "gitlab.niceprivate.com/golang/cook/util"
	"reflect"
)

type Ref_Model struct {
	R_Type reflect.Type
	PK     []*Ref_Field
	Fields []*Ref_Field
}

func NewRefModel(m interface{}) *Ref_Model {
	model := &Ref_Model{}
	model.R_Type = cook_util.REF_Must_struct_ptr(m)

	for i := 0; i < model.R_Type.NumField(); i++ {
		model.Fields = append(model.Fields, NewRefField(model, i))
	}

	return model
}
