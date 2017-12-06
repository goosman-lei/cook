package orm

import (
	"database/sql"
	"errors"
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

var (
	Err_invalid_model = errors.New("Invalid model type")
)

func NewRefModel(m interface{}) *Ref_Model {
	model := &Ref_Model{}
	if is_valid_model(m) {
		model.R_Type = reflect.TypeOf(m).Elem()
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

func is_valid_model(m interface{}) bool {
	r_model := reflect.TypeOf(m)
	is_struct_ptr := r_model.Kind() == reflect.Ptr && r_model.Elem().Kind() == reflect.Struct
	if !is_struct_ptr {
		return false
	}
	is_model := r_model.Elem().Implements(reflect.TypeOf((*Model)(nil)).Elem())
	if !is_model {
		return false
	}
	if m_field, has_field_m := r_model.Elem().FieldByName("M"); has_field_m {
		return m_field.Type == reflect.TypeOf((*M)(nil))
	}
	return false

}
func (g *God) is_model(model interface{}) bool {
	return reflect.TypeOf(model).Kind() == reflect.Ptr && reflect.TypeOf(model).Elem() == g.Model.R_Type
}

func (g *God) rows_to_model(rows *sql.Rows) (Model, error) {
	if rows == nil {
		return nil, nil
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, nil
	}

	datas, err := g.read_row(rows, cols)
	if err != nil {
		return nil, err
	}

	model := g.NewModel()
	r_model := reflect.Indirect(reflect.ValueOf(model))

	for col, val := range datas {
		if field, ok := g.Model.Mapping_with_column[col]; ok {
			cook_util.Assign_copy(r_model.FieldByName(field.R_StructField.Name), reflect.ValueOf(val))
		} else {
			if cook_util.Is_byte_slice(val) {
				model.Extra(col, string(val.([]byte)))
			} else {
				model.Extra(col, val)
			}
		}
	}

	return model, nil
}
func (g *God) rows_to_models(rows *sql.Rows) ([]Model, error) {
	if rows == nil {
		return nil, nil
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	models := []Model{}
	for {
		if !rows.Next() {
			break
		}

		datas, err := g.read_row(rows, cols)
		if err != nil {
			return nil, err
		}

		model := g.NewModel()
		r_model := reflect.Indirect(reflect.ValueOf(model))

		for col, val := range datas {
			if field, ok := g.Model.Mapping_with_column[col]; ok {
				cook_util.Assign_copy(r_model.FieldByName(field.R_StructField.Name), reflect.ValueOf(val))
			} else {
				if cook_util.Is_byte_slice(val) {
					model.Extra(col, string(val.([]byte)))
				} else {
					model.Extra(col, val)
				}
			}
		}
		models = append(models, model)
	}

	return models, nil

}

func (g *God) read_row(rows *sql.Rows, cols []string) (map[string]interface{}, error) {
	values := make([]interface{}, len(cols))
	valuesPtr := make([]interface{}, len(cols))
	for i, _ := range values {
		valuesPtr[i] = &values[i]
	}
	if err := rows.Scan(valuesPtr...); err != nil {
		return nil, err
	}
	mapping := make(map[string]interface{})
	for i, col := range cols {
		mapping[col] = *(valuesPtr[i].(*interface{}))
	}
	return mapping, nil
}
