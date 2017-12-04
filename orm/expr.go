package orm

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	err_skip_expr = errors.New("skip this expr")
)

type Expr struct {
	op      uint16
	field   string
	alias   string
	oprand1 interface{}
	oprand2 interface{}
	flags   uint64
}

type join_table struct {
	table     *Expr
	on        []*Expr
	join_type uint8
}

type SqlArgs []interface{}

func (s *SqlArgs) append(d ...interface{}) {
	*s = append(*s, d...)
}

const (
	// general expr
	OP_EQ uint16 = iota + 1
	OP_NE
	OP_GE
	OP_GT
	OP_LE
	OP_LT
	OP_BETWEEN
	OP_IS_NULL
	OP_NOT_NULL
	OP_LIKE
	OP_NOT_LIKE
	OP_IN
	OP_NOT_IN
	// logic
	OP_COMMA
	OP_AND
	OP_OR
	OP_NOT
	// set expr
	OP_ASSIGN
	OP_VALUES
	// literal
	OP_LITERAL
	OP_FIELD
	// from clause
	OP_TABLE
)

const (
	JOIN_TYPE_NORMAL uint8 = iota
	JOIN_TYPE_LEFT
	JOIN_TYPE_RIGHT
	JOIN_TYPE_INNER
	JOIN_TYPE_CROSS
)

const (
	FLAGS_ALL_SET  uint64 = 0XFFFFFFFFFFFFFFFF
	FLAGS_NONE_SET uint64 = 0x0

	// used 2 bits
	FLAGS_ORDERBY_MASK   = 0x3
	FLAGS_ORDERBY_OFFSET = 0
	FLAGS_ORDERBY_ASC    = 0x1 // DEFAULT
	FLAGS_ORDERBY_DESC   = 0x2
)

func multi_expr_to_string(exprs []*Expr, linker string, args *SqlArgs, buf *[]byte, off *int) error {
	is_first := true
	for _, expr := range exprs {
		if is_first {
			is_first = false
		} else {
			*off += copy((*buf)[*off:], linker)
		}
		if err := expr_to_string(expr, args, buf, off); err == err_skip_expr {
			*off -= len(linker)
		} else if err != nil {
			return err
		}
	}
	return nil
}
func expr_to_string(e *Expr, args *SqlArgs, buf *[]byte, off *int) error {
	switch e.op {
	case OP_EQ:
		if v, ok := e.oprand1.(*Expr); ok {
			*off += copy((*buf)[*off:], e.field+" = (")
			if err := expr_to_string(v, args, buf, off); err != nil {
				return err
			}
			(*buf)[*off] = ')'
			*off++
		} else {
			args.append(e.oprand1)
			*off += copy((*buf)[*off:], e.field+" = ?")
		}
	case OP_NE:
		if v, ok := e.oprand1.(*Expr); ok {
			*off += copy((*buf)[*off:], e.field+" != (")
			if err := expr_to_string(v, args, buf, off); err != nil {
				return err
			}
			(*buf)[*off] = ')'
			*off++
		} else {
			args.append(e.oprand1)
			*off += copy((*buf)[*off:], e.field+" != ?")
		}
	case OP_GE:
		if v, ok := e.oprand1.(*Expr); ok {
			*off += copy((*buf)[*off:], e.field+" >= (")
			if err := expr_to_string(v, args, buf, off); err != nil {
				return err
			}
			(*buf)[*off] = ')'
			*off++
		} else {
			args.append(e.oprand1)
			*off += copy((*buf)[*off:], e.field+" >= ?")
		}
	case OP_GT:
		if v, ok := e.oprand1.(*Expr); ok {
			*off += copy((*buf)[*off:], e.field+" > (")
			if err := expr_to_string(v, args, buf, off); err != nil {
				return err
			}
			(*buf)[*off] = ')'
			*off++
		} else {
			args.append(e.oprand1)
			*off += copy((*buf)[*off:], e.field+" > ?")
		}
	case OP_LE:
		if v, ok := e.oprand1.(*Expr); ok {
			*off += copy((*buf)[*off:], e.field+" <= (")
			if err := expr_to_string(v, args, buf, off); err != nil {
				return err
			}
			(*buf)[*off] = ')'
			*off++
		} else {
			args.append(e.oprand1)
			*off += copy((*buf)[*off:], e.field+" <= ?")
		}
	case OP_LT:
		if v, ok := e.oprand1.(*Expr); ok {
			*off += copy((*buf)[*off:], e.field+" < (")
			if err := expr_to_string(v, args, buf, off); err != nil {
				return err
			}
			(*buf)[*off] = ')'
			*off++
		} else {
			args.append(e.oprand1)
			*off += copy((*buf)[*off:], e.field+" < ?")
		}
	case OP_BETWEEN:
		if v, ok := e.oprand1.(*Expr); ok {
			*off += copy((*buf)[*off:], e.field+" BETWEEN (")
			if err := expr_to_string(v, args, buf, off); err != nil {
				return err
			}
			(*buf)[*off] = ')'
			*off++
		} else {
			args.append(e.oprand1)
			*off += copy((*buf)[*off:], e.field+" BETWEEN ?")
		}
		if v, ok := e.oprand2.(*Expr); ok {
			*off += copy((*buf)[*off:], e.field+" AND (")
			if err := expr_to_string(v, args, buf, off); err != nil {
				return err
			}
			(*buf)[*off] = ')'
			*off++
		} else {
			args.append(e.oprand2)
			*off += copy((*buf)[*off:], " AND ?")
		}
	case OP_IS_NULL:
		*off += copy((*buf)[*off:], e.field+" IS NULL")
	case OP_NOT_NULL:
		*off += copy((*buf)[*off:], e.field+" IS NOT NULL")
	case OP_LIKE:
		if v, ok := e.oprand1.(*Expr); ok {
			*off += copy((*buf)[*off:], e.field+" LIKE (")
			if err := expr_to_string(v, args, buf, off); err != nil {
				return err
			}
			(*buf)[*off] = ')'
			*off++
		} else {
			args.append(e.oprand1)
			*off += copy((*buf)[*off:], e.field+" LIKE ?")
		}
	case OP_NOT_LIKE:
		if v, ok := e.oprand1.(*Expr); ok {
			*off += copy((*buf)[*off:], e.field+" NOT LIKE (")
			if err := expr_to_string(v, args, buf, off); err != nil {
				return err
			}
			(*buf)[*off] = ')'
			*off++
		} else {
			args.append(e.oprand1)
			*off += copy((*buf)[*off:], e.field+" NOT LIKE ?")
		}
	case OP_IN:
		r_oprand1 := reflect.ValueOf(e.oprand1)
		switch r_oprand1.Kind() {
		case reflect.Slice, reflect.Array:
			if op_len := r_oprand1.Len(); op_len <= 0 {
				return err_skip_expr
			} else {
				r_args := reflect.ValueOf(args).Elem()
				*off += copy((*buf)[*off:], e.field+" IN(")
				for i, is_first := 0, true; i < op_len; i++ {
					r_args.Set(reflect.Append(r_args, r_oprand1.Index(i)))
					if is_first {
						is_first = false
					} else {
						(*buf)[*off] = ','
						*off++
						(*buf)[*off] = ' '
						*off++
					}
					(*buf)[*off] = '?'
					*off++
				}
				(*buf)[*off] = ')'
				*off++
			}
		default:
			return fmt.Errorf("<specify non-slice/non-array value to {%s in(?)}>", e.field)
		}
	case OP_NOT_IN:
		r_oprand1 := reflect.ValueOf(e.oprand1)
		switch r_oprand1.Kind() {
		case reflect.Slice, reflect.Array:
			if op_len := r_oprand1.Len(); op_len <= 0 {
				return err_skip_expr
			} else {
				r_args := reflect.ValueOf(args).Elem()
				*off += copy((*buf)[*off:], e.field+" NOT IN(")
				for i, is_first := 0, true; i < op_len; i++ {
					r_args.Set(reflect.Append(r_args, r_oprand1.Index(i)))
					if is_first {
						is_first = false
					} else {
						(*buf)[*off] = ','
						*off++
						(*buf)[*off] = ' '
						*off++
					}
					(*buf)[*off] = '?'
					*off++
				}
				(*buf)[*off] = ')'
				*off++
			}
		default:
			return fmt.Errorf("<specify non-slice/non-array value to {%s not in(?)}>", e.field)
		}
	case OP_COMMA:
		if v, ok := e.oprand1.([]*Expr); !ok {
			return fmt.Errorf("<specify non-[]*Expr value to COMMA expr>")
		} else {
			(*buf)[*off] = '('
			*off++
			if err := multi_expr_to_string(v, ", ", args, buf, off); err != nil {
				return err
			}
			(*buf)[*off] = ')'
			*off++
		}
	case OP_AND:
		if v, ok := e.oprand1.([]*Expr); !ok {
			return fmt.Errorf("<specify non-[]*Expr value to AND expr>")
		} else {
			(*buf)[*off] = '('
			*off++
			if err := multi_expr_to_string(v, " AND ", args, buf, off); err != nil {
				return err
			}
			(*buf)[*off] = ')'
			*off++
		}
	case OP_OR:
		if v, ok := e.oprand1.([]*Expr); !ok {
			return fmt.Errorf("<specify non-[]*Expr value to OR expr>")
		} else {
			(*buf)[*off] = '('
			*off++
			if err := multi_expr_to_string(v, " OR ", args, buf, off); err != nil {
				return err
			}
			(*buf)[*off] = ')'
			*off++
		}
	case OP_NOT:
		if v, ok := e.oprand1.(*Expr); !ok {
			return fmt.Errorf("<specify non-*Expr value to NOT expr>")
		} else {
			*off += copy((*buf)[*off:], "NOT (")
			if err := expr_to_string(v, args, buf, off); err != nil {
				return err
			}
			(*buf)[*off] = ')'
			*off++
		}
	case OP_ASSIGN:
		if v, ok := e.oprand1.(*Expr); ok {
			*off += copy((*buf)[*off:], e.field+" = (")
			if err := expr_to_string(v, args, buf, off); err != nil {
				return err
			}
			(*buf)[*off] = ')'
			*off++
		} else {
			args.append(e.oprand1)
			*off += copy((*buf)[*off:], e.field+" = ?")
		}
	case OP_VALUES:
		r_oprand1 := reflect.ValueOf(e.oprand1)
		switch r_oprand1.Kind() {
		case reflect.Slice, reflect.Array:
			if op_len := r_oprand1.Len(); op_len <= 0 {
				return err_skip_expr
			} else {
				r_args := reflect.ValueOf(args).Elem()
				(*buf)[*off] = '('
				*off++
				for i, is_first := 0, true; i < op_len; i++ {
					if is_first {
						is_first = false
					} else {
						(*buf)[*off] = ','
						*off++
						(*buf)[*off] = ' '
						*off++
					}
					if v, ok := r_oprand1.Index(i).Interface().(*Expr); ok {
						(*buf)[*off] = '('
						*off++
						if err := expr_to_string(v, args, buf, off); err != nil {
							return err
						}
						(*buf)[*off] = ')'
						*off++
					} else {
						r_args.Set(reflect.Append(r_args, r_oprand1.Index(i)))
						(*buf)[*off] = '?'
						*off++
					}
				}
				(*buf)[*off] = ')'
				*off++
			}
		default:
			return fmt.Errorf("<specify non-slice/non-array value to {%s VALUES(?)}>", e.field)
		}
	case OP_LITERAL:
		*off += copy((*buf)[*off:], e.oprand1.(string))
		if v, ok := e.oprand2.([]interface{}); ok && len(v) > 0 {
			args.append(v...)
		}
	case OP_FIELD:
		*off += copy((*buf)[*off:], e.field)
	case OP_TABLE:
		*off += copy((*buf)[*off:], e.field)
	default:
		return fmt.Errorf("<unknown expr type: %d>", e.op)
	}

	// process alias
	if len(e.alias) > 0 {
		*off += copy((*buf)[*off:], " AS "+e.alias)
	}

	// process table join
	if e.op == OP_TABLE && e.oprand1 != nil {
		for _, j_table := range e.oprand1.([]*join_table) {
			var join_word string = " JOIN "
			switch j_table.join_type {
			case JOIN_TYPE_LEFT:
				join_word = " LEFT JOIN "
			case JOIN_TYPE_RIGHT:
				join_word = " RIGHT JOIN "
			case JOIN_TYPE_INNER:
				join_word = " INNER JOIN "
			case JOIN_TYPE_CROSS:
				join_word = " CROSS JOIN "
			}

			*off += copy((*buf)[*off:], join_word)
			if err := expr_to_string(j_table.table, args, buf, off); err != nil {
				return err
			}
			if len(j_table.on) > 0 {
				*off += copy((*buf)[*off:], " ON ")
				if err := multi_expr_to_string(j_table.on, " AND ", args, buf, off); err != nil {
					return err
				}
			}
		}
	}

	// process flags
	switch e.get_flag(FLAGS_ORDERBY_MASK, FLAGS_ORDERBY_OFFSET) {
	case FLAGS_ORDERBY_ASC:
		*off += copy((*buf)[*off:], " ASC")
	case FLAGS_ORDERBY_DESC:
		*off += copy((*buf)[*off:], " DESC")
	}

	return nil
}

func (e *Expr) Alias(alias string) *Expr {
	e.alias = alias
	return e
}

func (e *Expr) set_flag(val, mask uint64, offset uint) {
	e.flags = ((FLAGS_ALL_SET ^ (mask << offset)) & e.flags) | (val & mask << offset)
}

func (e *Expr) get_flag(mask uint64, offset uint) uint64 {
	return (FLAGS_NONE_SET ^ (mask << offset)) & e.flags >> offset
}

func (e *Expr) Asc() *Expr {
	e.set_flag(FLAGS_ORDERBY_ASC, FLAGS_ORDERBY_MASK, FLAGS_ORDERBY_OFFSET)
	return e
}

func (e *Expr) Desc() *Expr {
	e.set_flag(FLAGS_ORDERBY_DESC, FLAGS_ORDERBY_MASK, FLAGS_ORDERBY_OFFSET)
	return e
}

func E_eq(f string, d interface{}) *Expr {
	return &Expr{op: OP_EQ, field: f, oprand1: d}
}

func E_ne(f string, d interface{}) *Expr {
	return &Expr{op: OP_NE, field: f, oprand1: d}
}

func E_ge(f string, d interface{}) *Expr {
	return &Expr{op: OP_GE, field: f, oprand1: d}
}

func E_gt(f string, d interface{}) *Expr {
	return &Expr{op: OP_GT, field: f, oprand1: d}
}

func E_le(f string, d interface{}) *Expr {
	return &Expr{op: OP_LE, field: f, oprand1: d}
}

func E_lt(f string, d interface{}) *Expr {
	return &Expr{op: OP_LT, field: f, oprand1: d}
}

func E_between(f string, d1, d2 interface{}) *Expr {
	return &Expr{op: OP_BETWEEN, field: f, oprand1: d1, oprand2: d2}
}

func E_is_null(f string) *Expr {
	return &Expr{op: OP_IS_NULL, field: f}
}

func E_not_null(f string) *Expr {
	return &Expr{op: OP_NOT_NULL, field: f}
}

func E_like(f string, d string) *Expr {
	return &Expr{op: OP_LIKE, field: f, oprand1: d}
}

func E_not_like(f string, d string) *Expr {
	return &Expr{op: OP_NOT_LIKE, field: f, oprand1: d}
}

func E_in(f string, d interface{}) *Expr {
	return &Expr{op: OP_IN, field: f, oprand1: d}
}

func E_not_in(f string, d interface{}) *Expr {
	return &Expr{op: OP_NOT_IN, field: f, oprand1: d}
}

func E_comma(d ...*Expr) *Expr {
	return &Expr{op: OP_COMMA, oprand1: d}
}

func E_and(d ...*Expr) *Expr {
	return &Expr{op: OP_AND, oprand1: d}
}

func E_or(d ...*Expr) *Expr {
	return &Expr{op: OP_OR, oprand1: d}
}

func E_not(d *Expr) *Expr {
	return &Expr{op: OP_NOT, oprand1: d}
}

func E_assign(f string, d interface{}) *Expr {
	return &Expr{op: OP_ASSIGN, field: f, oprand1: d}
}

func E_values(d ...interface{}) *Expr {
	return &Expr{op: OP_VALUES, oprand1: d}
}

func E_field(f string) *Expr {
	return &Expr{op: OP_FIELD, field: f}
}

func E_asc(f string) *Expr {
	return E_field(f).Asc()
}

func E_desc(f string) *Expr {
	return E_field(f).Desc()
}

func E_fields(fs ...string) []*Expr {
	es := []*Expr{}
	for _, f := range fs {
		es = append(es, &Expr{op: OP_FIELD, field: f})
	}
	return es
}

func E_literal(d string, args ...interface{}) *Expr {
	return &Expr{op: OP_LITERAL, oprand1: d, oprand2: args}
}

func E_table(t string) *Expr {
	return &Expr{op: OP_TABLE, field: t}
}

func (e *Expr) join(join_type uint8, t *Expr, on ...*Expr) *Expr {
	if e.oprand1 == nil {
		e.oprand1 = make([]*join_table, 0)
	}
	e.oprand1 = append(e.oprand1.([]*join_table), &join_table{
		table:     t,
		on:        on,
		join_type: join_type,
	})
	return e

}

func (e *Expr) Join(t *Expr, on ...*Expr) *Expr {
	return e.join(JOIN_TYPE_NORMAL, t, on...)
}

func (e *Expr) LeftJoin(t *Expr, on ...*Expr) *Expr {
	return e.join(JOIN_TYPE_LEFT, t, on...)
}

func (e *Expr) RightJoin(t *Expr, on ...*Expr) *Expr {
	return e.join(JOIN_TYPE_RIGHT, t, on...)
}

func (e *Expr) InnerJoin(t *Expr, on ...*Expr) *Expr {
	return e.join(JOIN_TYPE_INNER, t, on...)
}

func (e *Expr) CrossJoin(t *Expr, on ...*Expr) *Expr {
	return e.join(JOIN_TYPE_CROSS, t, on...)
}
