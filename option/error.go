package option

import (
	"fmt"
)

type OptionError struct {
	ErrType uint8
	ErrMsg  string
}

const (
	Err_Target_Not_StructPtr uint8 = iota + 1
	Err_Field_Invalid
	Err_Field_CanNot_Set
	Err_Type_Incompatible
)

var (
	errors_mapping map[uint8]string = map[uint8]string{
		Err_Target_Not_StructPtr: "target options is not pointer to struct",
		Err_Field_Invalid:        "Field %s is invalid",
		Err_Field_CanNot_Set:     "Field %s can't be settable",
		Err_Type_Incompatible:    "Type Incompatible for field[%s]: %s",
	}
)

func (o OptionError) Error() string {
	return o.ErrMsg
}

func NewError(errtype uint8, argv ...interface{}) OptionError {
	return OptionError{
		ErrType: errtype,
		ErrMsg:  fmt.Sprintf(errors_mapping[errtype], argv...),
	}
}
