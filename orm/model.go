package orm

import (
	cook_util "gitlab.niceprivate.com/golang/cook/util"
	"reflect"
)

type Model interface {
	InitExtra()
	Extra(k string, v interface{})
	Extras(map[string]interface{})
	Interface(k string) interface{}
	Bool(k string) bool
	Int(k string) int
	Int8(k string) int8
	Int16(k string) int16
	Int32(k string) int32
	Int64(k string) int64
	Uint(k string) uint
	Uint8(k string) uint8
	Uint16(k string) uint16
	Uint32(k string) uint32
	Uint64(k string) uint64
	Float32(k string) float32
	Float64(k string) float64
	String(k string) string
}

type M struct {
	E map[string]interface{}
}

func (g *God) NewModel() Model {
	model := g.Factory()
	field := reflect.Indirect(reflect.ValueOf(model)).FieldByName("M")
	if field.IsNil() {
		m := new(M)
		m.InitExtra()
		field.Set(reflect.ValueOf(m))
	}
	return model
}

func (m *M) InitExtra() {
	m.E = make(map[string]interface{})
}

func (m *M) Extra(k string, v interface{}) {
	m.E[k] = v
}
func (m *M) Extras(datas map[string]interface{}) {
	for k, v := range datas {
		m.E[k] = v
	}
}
func (m *M) Exists(k string) bool {
	_, ok := m.E[k]
	return ok
}
func (m *M) Interface(k string) interface{} {
	if v, ok := m.E[k]; ok {
		return v
	} else {
		return nil
	}
}
func (m *M) Bool(k string) bool {
	return cook_util.As_bool(m.E[k])
}
func (m *M) Int(k string) int {
	return cook_util.As_int(m.E[k])
}
func (m *M) Int8(k string) int8 {
	return cook_util.As_int8(m.E[k])
}
func (m *M) Int16(k string) int16 {
	return cook_util.As_int16(m.E[k])
}
func (m *M) Int32(k string) int32 {
	return cook_util.As_int32(m.E[k])
}
func (m *M) Int64(k string) int64 {
	return cook_util.As_int64(m.E[k])
}
func (m *M) Uint(k string) uint {
	return cook_util.As_uint(m.E[k])
}
func (m *M) Uint8(k string) uint8 {
	return cook_util.As_uint8(m.E[k])
}
func (m *M) Uint16(k string) uint16 {
	return cook_util.As_uint16(m.E[k])
}
func (m *M) Uint32(k string) uint32 {
	return cook_util.As_uint32(m.E[k])
}
func (m *M) Uint64(k string) uint64 {
	return cook_util.As_uint64(m.E[k])
}
func (m *M) Float32(k string) float32 {
	return cook_util.As_float32(m.E[k])
}
func (m *M) Float64(k string) float64 {
	return cook_util.As_float64(m.E[k])
}
func (m *M) String(k string) string {
	return cook_util.As_string(m.E[k])
}
