package option

import (
	"reflect"
)

type I_Options interface{}

type Option func(o I_Options) error

func Apply(target I_Options, opts ...Option) error {
	for _, opt := range opts {
		if err := opt(target); err != nil {
			return err
		}
	}
	return nil
}

func Opt(name string, value interface{}) Option {
	return func(o I_Options) error {
		return set_field(o, name, reflect.ValueOf(value))
	}
}

func set_field(target interface{}, name string, value reflect.Value) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = NewError(Err_Type_Incompatible, name, r)
		}
	}()
	var (
		r_ptr, r_struct, r_field reflect.Value
	)
	if r_ptr = reflect.ValueOf(target); r_ptr.Kind() != reflect.Ptr {
		return NewError(Err_Target_Not_StructPtr)
	}
	if r_struct = r_ptr.Elem(); r_struct.Kind() != reflect.Struct {
		return NewError(Err_Target_Not_StructPtr)
	}
	if r_field = r_struct.FieldByName(name); !r_field.IsValid() {
		return NewError(Err_Field_Invalid, name)
	}
	if !r_field.CanSet() {
		return NewError(Err_Field_CanNot_Set, name)
	}

	r_field.Set(value)
	return nil
}
