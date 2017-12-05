package orm

import (
	"strconv"
)

type Model interface {
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
	if e, ok := m.E[k]; !ok {
		return false
	} else {
		switch v := e.(type) {
		case bool:
			return v
		case int:
			return v != 0
		case int8:
			return v != 0
		case int16:
			return v != 0
		case int32:
			return v != 0
		case int64:
			return v != 0
		case uint:
			return v != 0
		case uint8:
			return v != 0
		case uint16:
			return v != 0
		case uint32:
			return v != 0
		case uint64:
			return v != 0
		case float32:
			return v != 0.0
		case float64:
			return v != 0.0
		case string:
			return len(v) > 0
		default:
			return false
		}
	}
}
func (m *M) Int(k string) int {
	if e, ok := m.E[k]; !ok {
		return 0
	} else {
		switch v := e.(type) {
		case bool:
			if v {
				return 1
			} else {
				return 0
			}
		case int:
			return v
		case int8:
			return int(v)
		case int16:
			return int(v)
		case int32:
			return int(v)
		case int64:
			return int(v)
		case uint:
			return int(v)
		case uint8:
			return int(v)
		case uint16:
			return int(v)
		case uint32:
			return int(v)
		case uint64:
			return int(v)
		case float32:
			return int(v)
		case float64:
			return int(v)
		case string:
			if iv, err := strconv.ParseInt(v, 10, 64); err != nil {
				return 0
			} else {
				return int(iv)
			}
		default:
			return 0
		}
	}
}
func (m *M) Int8(k string) int8 {
	if e, ok := m.E[k]; !ok {
		return 0
	} else {
		switch v := e.(type) {
		case bool:
			if v {
				return 1
			} else {
				return 0
			}
		case int:
			return int8(v)
		case int8:
			return v
		case int16:
			return int8(v)
		case int32:
			return int8(v)
		case int64:
			return int8(v)
		case uint:
			return int8(v)
		case uint8:
			return int8(v)
		case uint16:
			return int8(v)
		case uint32:
			return int8(v)
		case uint64:
			return int8(v)
		case float32:
			return int8(v)
		case float64:
			return int8(v)
		case string:
			if iv, err := strconv.ParseInt(v, 10, 64); err != nil {
				return 0
			} else {
				return int8(iv)
			}
		default:
			return 0
		}
	}
}
func (m *M) Int16(k string) int16 {
	if e, ok := m.E[k]; !ok {
		return 0
	} else {
		switch v := e.(type) {
		case bool:
			if v {
				return 1
			} else {
				return 0
			}
		case int:
			return int16(v)
		case int8:
			return int16(v)
		case int16:
			return v
		case int32:
			return int16(v)
		case int64:
			return int16(v)
		case uint:
			return int16(v)
		case uint8:
			return int16(v)
		case uint16:
			return int16(v)
		case uint32:
			return int16(v)
		case uint64:
			return int16(v)
		case float32:
			return int16(v)
		case float64:
			return int16(v)
		case string:
			if iv, err := strconv.ParseInt(v, 10, 64); err != nil {
				return 0
			} else {
				return int16(iv)
			}
		default:
			return 0
		}
	}
}
func (m *M) Int32(k string) int32 {
	if e, ok := m.E[k]; !ok {
		return 0
	} else {
		switch v := e.(type) {
		case bool:
			if v {
				return 1
			} else {
				return 0
			}
		case int:
			return int32(v)
		case int8:
			return int32(v)
		case int16:
			return int32(v)
		case int32:
			return v
		case int64:
			return int32(v)
		case uint:
			return int32(v)
		case uint8:
			return int32(v)
		case uint16:
			return int32(v)
		case uint32:
			return int32(v)
		case uint64:
			return int32(v)
		case float32:
			return int32(v)
		case float64:
			return int32(v)
		case string:
			if iv, err := strconv.ParseInt(v, 10, 64); err != nil {
				return 0
			} else {
				return int32(iv)
			}
		default:
			return 0
		}
	}
}
func (m *M) Int64(k string) int64 {
	if e, ok := m.E[k]; !ok {
		return 0
	} else {
		switch v := e.(type) {
		case bool:
			if v {
				return 1
			} else {
				return 0
			}
		case int:
			return int64(v)
		case int8:
			return int64(v)
		case int16:
			return int64(v)
		case int32:
			return int64(v)
		case int64:
			return v
		case uint:
			return int64(v)
		case uint8:
			return int64(v)
		case uint16:
			return int64(v)
		case uint32:
			return int64(v)
		case uint64:
			return int64(v)
		case float32:
			return int64(v)
		case float64:
			return int64(v)
		case string:
			if iv, err := strconv.ParseInt(v, 10, 64); err != nil {
				return 0
			} else {
				return int64(iv)
			}
		default:
			return 0
		}
	}
}
func (m *M) Uint(k string) uint {
	if e, ok := m.E[k]; !ok {
		return 0
	} else {
		switch v := e.(type) {
		case bool:
			if v {
				return 1
			} else {
				return 0
			}
		case int:
			return uint(v)
		case int8:
			return uint(v)
		case int16:
			return uint(v)
		case int32:
			return uint(v)
		case int64:
			return uint(v)
		case uint:
			return v
		case uint8:
			return uint(v)
		case uint16:
			return uint(v)
		case uint32:
			return uint(v)
		case uint64:
			return uint(v)
		case float32:
			return uint(v)
		case float64:
			return uint(v)
		case string:
			if iv, err := strconv.ParseInt(v, 10, 64); err != nil {
				return 0
			} else {
				return uint(iv)
			}
		default:
			return 0
		}
	}
}
func (m *M) Uint8(k string) uint8 {
	if e, ok := m.E[k]; !ok {
		return 0
	} else {
		switch v := e.(type) {
		case bool:
			if v {
				return 1
			} else {
				return 0
			}
		case int:
			return uint8(v)
		case int8:
			return uint8(v)
		case int16:
			return uint8(v)
		case int32:
			return uint8(v)
		case int64:
			return uint8(v)
		case uint:
			return uint8(v)
		case uint8:
			return v
		case uint16:
			return uint8(v)
		case uint32:
			return uint8(v)
		case uint64:
			return uint8(v)
		case float32:
			return uint8(v)
		case float64:
			return uint8(v)
		case string:
			if iv, err := strconv.ParseInt(v, 10, 64); err != nil {
				return 0
			} else {
				return uint8(iv)
			}
		default:
			return 0
		}
	}
}
func (m *M) Uint16(k string) uint16 {
	if e, ok := m.E[k]; !ok {
		return 0
	} else {
		switch v := e.(type) {
		case bool:
			if v {
				return 1
			} else {
				return 0
			}
		case int:
			return uint16(v)
		case int8:
			return uint16(v)
		case int16:
			return uint16(v)
		case int32:
			return uint16(v)
		case int64:
			return uint16(v)
		case uint:
			return uint16(v)
		case uint8:
			return uint16(v)
		case uint16:
			return v
		case uint32:
			return uint16(v)
		case uint64:
			return uint16(v)
		case float32:
			return uint16(v)
		case float64:
			return uint16(v)
		case string:
			if iv, err := strconv.ParseInt(v, 10, 64); err != nil {
				return 0
			} else {
				return uint16(iv)
			}
		default:
			return 0
		}
	}
}
func (m *M) Uint32(k string) uint32 {
	if e, ok := m.E[k]; !ok {
		return 0
	} else {
		switch v := e.(type) {
		case bool:
			if v {
				return 1
			} else {
				return 0
			}
		case int:
			return uint32(v)
		case int8:
			return uint32(v)
		case int16:
			return uint32(v)
		case int32:
			return uint32(v)
		case int64:
			return uint32(v)
		case uint:
			return uint32(v)
		case uint8:
			return uint32(v)
		case uint16:
			return uint32(v)
		case uint32:
			return uint32(v)
		case uint64:
			return uint32(v)
		case float32:
			return uint32(v)
		case float64:
			return uint32(v)
		case string:
			if iv, err := strconv.ParseInt(v, 10, 64); err != nil {
				return 0
			} else {
				return uint32(iv)
			}
		default:
			return 0
		}
	}
}
func (m *M) Uint64(k string) uint64 {
	if e, ok := m.E[k]; !ok {
		return 0
	} else {
		switch v := e.(type) {
		case bool:
			if v {
				return 1
			} else {
				return 0
			}
		case int:
			return uint64(v)
		case int8:
			return uint64(v)
		case int16:
			return uint64(v)
		case int32:
			return uint64(v)
		case int64:
			return uint64(v)
		case uint:
			return uint64(v)
		case uint8:
			return uint64(v)
		case uint16:
			return uint64(v)
		case uint32:
			return uint64(v)
		case uint64:
			return uint64(v)
		case float32:
			return uint64(v)
		case float64:
			return uint64(v)
		case string:
			if iv, err := strconv.ParseInt(v, 10, 64); err != nil {
				return 0
			} else {
				return uint64(iv)
			}
		default:
			return 0
		}
	}
}
func (m *M) Float32(k string) float32 {
	if e, ok := m.E[k]; !ok {
		return 0.0
	} else {
		switch v := e.(type) {
		case bool:
			if v {
				return 1.0
			} else {
				return 0.0
			}
		case int:
			return float32(v)
		case int8:
			return float32(v)
		case int16:
			return float32(v)
		case int32:
			return float32(v)
		case int64:
			return float32(v)
		case uint:
			return float32(v)
		case uint8:
			return float32(v)
		case uint16:
			return float32(v)
		case uint32:
			return float32(v)
		case uint64:
			return float32(v)
		case float32:
			return float32(v)
		case float64:
			return float32(v)
		case string:
			if fv, err := strconv.ParseFloat(v, 32); err != nil {
				return 0.0
			} else {
				return float32(fv)
			}
		default:
			return 0.0
		}
	}
}
func (m *M) Float64(k string) float64 {
	if e, ok := m.E[k]; !ok {
		return 0.0
	} else {
		switch v := e.(type) {
		case bool:
			if v {
				return 1.0
			} else {
				return 0.0
			}
		case int:
			return float64(v)
		case int8:
			return float64(v)
		case int16:
			return float64(v)
		case int32:
			return float64(v)
		case int64:
			return float64(v)
		case uint:
			return float64(v)
		case uint8:
			return float64(v)
		case uint16:
			return float64(v)
		case uint32:
			return float64(v)
		case uint64:
			return float64(v)
		case float32:
			return float64(v)
		case float64:
			return float64(v)
		case string:
			if fv, err := strconv.ParseFloat(v, 64); err != nil {
				return 0.0
			} else {
				return float64(fv)
			}
		default:
			return 0.0
		}
	}
}
func (m *M) String(k string) string {
	if e, ok := m.E[k]; !ok {
		return ""
	} else {
		switch v := e.(type) {
		case bool:
			return strconv.FormatBool(v)
		case int:
			return strconv.FormatInt(int64(v), 10)
		case int8:
			return strconv.FormatInt(int64(v), 10)
		case int16:
			return strconv.FormatInt(int64(v), 10)
		case int32:
			return strconv.FormatInt(int64(v), 10)
		case int64:
			return strconv.FormatInt(v, 10)
		case uint:
			return strconv.FormatUint(uint64(v), 10)
		case uint8:
			return strconv.FormatUint(uint64(v), 10)
		case uint16:
			return strconv.FormatUint(uint64(v), 10)
		case uint32:
			return strconv.FormatUint(uint64(v), 10)
		case uint64:
			return strconv.FormatUint(v, 10)
		case float32:
			return strconv.FormatFloat(float64(v), 'f', -1, 32)
		case float64:
			return strconv.FormatFloat(v, 'f', -1, 64)
		case string:
			return v
		default:
			return ""
		}
	}
}
