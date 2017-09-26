package connector

import (
	"bytes"
	"fmt"
	"strings"
)

type MArgs struct {
	args []interface{}
}

func (ma *MArgs) append(args ...interface{}) {
	ma.args = append(ma.args, args...)
}

type MTable struct {
	join          int
	name          string
	alias         string
	idx_hint_list []string
	on            MExpr
}

type MExpr struct {
	operator uint8
	field    string
	oprand1  interface{}
	oprand2  interface{}
	flags    uint64
}

const (
	MYSQL_OP_NONE = iota
	MYSQL_OP_EQ
	MYSQL_OP_NE
	MYSQL_OP_GE
	MYSQL_OP_GT
	MYSQL_OP_LE
	MYSQL_OP_LT
	MYSQL_OP_BETWEEN
	MYSQL_OP_IS_NULL
	MYSQL_OP_NOT_NULL
	MYSQL_OP_LIKE
	MYSQL_OP_NOT_LIKE
	MYSQL_OP_IN
	MYSQL_OP_NOT_IN
	MYSQL_OP_LITERAL
	MYSQL_OP_AND
	MYSQL_OP_OR
	MYSQL_OP_NOT
	MYSQL_OP_FIELD
	MYSQL_OP_ALIAS
	MYSQL_OP_EXPR_ALIAS
	MYSQL_OP_ASSIGN
	MYSQL_OP_ASSIGN_VALUE
)

const (
	MYSQL_LOGIC_AND = iota
	MYSQL_LOGIC_OR
	MYSQL_LOGIC_COMMA
)

const (
	MYSQL_JOIN_NONE = iota
	MYSQL_JOIN_NORMAL
	MYSQL_JOIN_INNER
	MYSQL_JOIN_CROSS
	MYSQL_JOIN_LEFT
	MYSQL_JOIN_RIGHT
)

const (
	MYSQL_FLAGS_ORDER_OFFSET = 0
	MYSQL_FLAGS_ORDER_MASK   = 0x3
	MYSQL_FLAGS_ORDER_UNMASK = 0xFFFFFFFFFFFFFFFC
	MYSQL_FLAGS_ORDER_NONE   = 0x0
	MYSQL_FLAGS_ORDER_ASC    = 0x1
	MYSQL_FLAGS_ORDER_DESC   = 0x2
)

var EmptyMExpr MExpr
var EmptyMExprSlice []MExpr

func parseDelete(table string, where []MExpr, orderBy []MExpr, limit int) (query string, args []interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	ma := &MArgs{}
	buf := new(bytes.Buffer)

	parseMTables(MT_Slice_Name(table), "DELETE FROM", buf, ma)
	parseMExprs(MYSQL_LOGIC_AND, where, " WHERE ", buf, ma)
	parseMExprs(MYSQL_LOGIC_COMMA, orderBy, " ORDER BY ", buf, ma)
	parseLimit(limit, -1, buf)

	query = buf.String()
	args = ma.args
	return

}

func parseInsert(table string, sets []MExpr, onDup []MExpr) (query string, args []interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	ma := &MArgs{}
	buf := new(bytes.Buffer)

	parseMTables(MT_Slice_Name(table), "INSERT INTO", buf, ma)
	parseMExprs(MYSQL_LOGIC_COMMA, sets, " SET ", buf, ma)
	parseMExprs(MYSQL_LOGIC_COMMA, onDup, " ON DUPLICATE KEY UPDATE ", buf, ma)

	query = buf.String()
	args = ma.args
	return
}

func parseUpdate(tables []MTable, sets []MExpr, where []MExpr, orderBy []MExpr, limit int) (query string, args []interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	ma := &MArgs{}
	buf := new(bytes.Buffer)

	parseMTables(tables, "UPDATE", buf, ma)
	parseMExprs(MYSQL_LOGIC_COMMA, sets, " SET ", buf, ma)
	parseMExprs(MYSQL_LOGIC_AND, where, " WHERE ", buf, ma)
	parseMExprs(MYSQL_LOGIC_COMMA, orderBy, " ORDER BY ", buf, ma)
	parseLimit(limit, -1, buf)

	query = buf.String()
	args = ma.args
	return
}

func parseQuery(tables []MTable, fields []MExpr, where []MExpr, limit, offset int, orderBy, groupBy, having []MExpr) (query string, args []interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	ma := &MArgs{}
	buf := new(bytes.Buffer)

	parseMExprs(MYSQL_LOGIC_COMMA, fields, "SELECT ", buf, ma)
	parseMTables(tables, " FROM", buf, ma)
	parseMExprs(MYSQL_LOGIC_AND, where, " WHERE ", buf, ma)
	parseMExprs(MYSQL_LOGIC_COMMA, groupBy, " GROUP BY ", buf, ma)
	parseMExprs(MYSQL_LOGIC_AND, having, " HAVING ", buf, ma)
	parseMExprs(MYSQL_LOGIC_COMMA, orderBy, " ORDER BY ", buf, ma)
	parseLimit(limit, offset, buf)

	query = buf.String()
	args = ma.args
	return
}

func parseMTables(tables []MTable, prefix string, buf *bytes.Buffer, ma *MArgs) {
	isFirst := true

	for _, table := range tables {
		if isFirst {
			fmt.Fprintf(buf, "%s", prefix)
			isFirst = false
		}

		switch table.join {
		case MYSQL_JOIN_NORMAL:
			fmt.Fprintf(buf, " JOIN")
		case MYSQL_JOIN_INNER:
			fmt.Fprintf(buf, " INNER JOIN")
		case MYSQL_JOIN_CROSS:
			fmt.Fprintf(buf, " CROSS JOIN")
		case MYSQL_JOIN_LEFT:
			fmt.Fprintf(buf, " LEFT JOIN")
		case MYSQL_JOIN_RIGHT:
			fmt.Fprintf(buf, " RIGHT JOIN")
		}

		fmt.Fprintf(buf, " %s", table.name)

		if len(table.alias) > 0 {
			fmt.Fprintf(buf, " AS %s", table.alias)
		}

		if len(table.idx_hint_list) > 0 {
			fmt.Fprintf(buf, " %s", strings.Join(table.idx_hint_list, ", "))
		}

		if table.on.operator != MYSQL_OP_NONE {
			parseMExprWithPrefix(table.on, " ON ", buf, ma)
		}
	}
}

func parseLimit(limit, offset int, buf *bytes.Buffer) {
	if limit > 0 {
		fmt.Fprintf(buf, " LIMIT %d", limit)
		if offset > 0 {
			fmt.Fprintf(buf, " OFFSET %d", offset)
		}
	}
}

func parseMExprs(logic int, conds []MExpr, prefix string, buf *bytes.Buffer, ma *MArgs) {
	var (
		isFirst    bool = true
		joinPrefix string
	)

	if logic == MYSQL_LOGIC_AND {
		joinPrefix = " AND "
	} else if logic == MYSQL_LOGIC_OR {
		joinPrefix = " OR "
	} else {
		joinPrefix = ", "
	}

	for _, cond := range conds {
		if isFirst {
			parseMExprWithPrefix(cond, prefix, buf, ma)
			isFirst = false
		} else {
			parseMExprWithPrefix(cond, joinPrefix, buf, ma)
		}
	}
}

func parseMExpr(me MExpr, buf *bytes.Buffer, ma *MArgs) {
	parseMExprWithPrefix(me, "", buf, ma)
}

func parseMExprWithPrefix(me MExpr, prefix string, buf *bytes.Buffer, ma *MArgs) {
	switch me.operator {
	case MYSQL_OP_EQ:
		ma.append(me.oprand1)
		fmt.Fprintf(buf, "%s%s = ?", prefix, me.field)
	case MYSQL_OP_NE:
		ma.append(me.oprand1)
		fmt.Fprintf(buf, "%s%s != ?", prefix, me.field)
	case MYSQL_OP_GE:
		ma.append(me.oprand1)
		fmt.Fprintf(buf, "%s%s >= ?", prefix, me.field)
	case MYSQL_OP_GT:
		ma.append(me.oprand1)
		fmt.Fprintf(buf, "%s%s > ?", prefix, me.field)
	case MYSQL_OP_LE:
		ma.append(me.oprand1)
		fmt.Fprintf(buf, "%s%s <= ?", prefix, me.field)
	case MYSQL_OP_LT:
		ma.append(me.oprand1)
		fmt.Fprintf(buf, "%s%s < ?", prefix, me.field)
	case MYSQL_OP_BETWEEN:
		ma.append(me.oprand1, me.oprand2)
		fmt.Fprintf(buf, "%s%s BETWEEN ? AND ?", prefix, me.field)
	case MYSQL_OP_IS_NULL:
		fmt.Fprintf(buf, "%s%s IS NULL", prefix, me.field)
	case MYSQL_OP_NOT_NULL:
		fmt.Fprintf(buf, "%s%s IS NOT NULL", prefix, me.field)
	case MYSQL_OP_LIKE:
		ma.append(me.oprand1)
		fmt.Fprintf(buf, "%s%s LIKE ?", prefix, me.field)
	case MYSQL_OP_NOT_LIKE:
		ma.append(me.oprand1)
		fmt.Fprintf(buf, "%s%s NOT LIKE ?", prefix, me.field)
	case MYSQL_OP_IN:
		if v, ok := me.oprand1.([]interface{}); ok {
			if len(v) > 0 {
				ma.append(v...)
				fmt.Fprintf(buf, "%s%s IN(%s)", prefix, me.field, strings.Join(strings.Split(strings.Repeat("?", len(v)), ""), ","))
			}
		} else {
			panic(fmt.Errorf("mysql_op_in expetected arg is slice, but: %q", v))
		}
	case MYSQL_OP_NOT_IN:
		if v, ok := me.oprand1.([]interface{}); ok {
			if len(v) > 0 {
				ma.append(v...)
				fmt.Fprintf(buf, "%s%s NOT IN(%s)", prefix, me.field, strings.Join(strings.Split(strings.Repeat("?", len(v)), ""), ","))
			}
		} else {
			panic(fmt.Errorf("mysql_op_not_in expetected arg is slice, but: %q", v))
		}
	// literal string must not contain "?"
	case MYSQL_OP_LITERAL:
		fmt.Fprintf(buf, "%s%s", prefix, me.oprand1)
	case MYSQL_OP_AND:
		if v, ok := me.oprand1.([]MExpr); ok {
			fmt.Fprintf(buf, "%s(", prefix)
			parseMExprs(MYSQL_LOGIC_AND, v, "", buf, ma)
			fmt.Fprintf(buf, ")")
		} else {
			panic(fmt.Errorf("mysql_op_and expetected arg is []MExpr, but: %q", v))
		}
	case MYSQL_OP_OR:
		if v, ok := me.oprand1.([]MExpr); ok {
			fmt.Fprintf(buf, "%s(", prefix)
			parseMExprs(MYSQL_LOGIC_OR, v, "", buf, ma)
			fmt.Fprintf(buf, ")")
		} else {
			panic(fmt.Errorf("mysql_op_or expetected arg is []MExpr, but: %q", v))
		}
	case MYSQL_OP_NOT:
		if v, ok := me.oprand1.(MExpr); ok {
			fmt.Fprintf(buf, "%sNOT (", prefix)
			parseMExpr(v, buf, ma)
			fmt.Fprintf(buf, ")")
		} else {
			panic(fmt.Errorf("mysql_op_not expetected arg is MExpr, but: %q", v))
		}
	case MYSQL_OP_FIELD:
		fmt.Fprintf(buf, "%s%s", prefix, me.field)
	case MYSQL_OP_ALIAS:
		fmt.Fprintf(buf, "%s%s AS %s", prefix, me.field, me.oprand1.(string))
	case MYSQL_OP_EXPR_ALIAS:
		if v, ok := me.oprand1.(MExpr); ok {
			fmt.Fprintf(buf, "%s(", prefix)
			parseMExpr(v, buf, ma)
			fmt.Fprintf(buf, ") AS %s", me.oprand2.(string))
		} else {
			panic(fmt.Errorf("mysql_op_expr_alias expetected arg is MExpr, but: %q", v))
		}
	case MYSQL_OP_ASSIGN:
		if v, ok := me.oprand1.(MExpr); ok {
			fmt.Fprintf(buf, "%s%s = (", prefix, me.field)
			parseMExpr(v, buf, ma)
			fmt.Fprintf(buf, ")")
		} else {
			panic(fmt.Errorf("mysql_op_assign expetected arg is MExpr, but: %q", v))
		}
	case MYSQL_OP_ASSIGN_VALUE:
		fmt.Fprintf(buf, "%s%s = ?", prefix, me.field)
		ma.append(me.oprand1)
	default:
		panic(fmt.Errorf("unknown expr: %q", me))
	}
	if buf.Len() > 0 {
		switch me.FlagsOrder() {
		case MYSQL_FLAGS_ORDER_ASC:
			fmt.Fprintf(buf, " ASC")
		case MYSQL_FLAGS_ORDER_DESC:
			fmt.Fprintf(buf, " DESC")
		}
	}
}

func MT_Slice_Names(names []string) (r []MTable) {
	for _, name := range names {
		r = append(r, MT_Name(name))
	}
	return
}

func MT_Slice_Name(name string) (r []MTable) {
	r = append(r, MT_Name(name))
	return
}

func MT_Name(name string) MTable {
	return MTable{name: name}
}

func ME_Slice_Fields(fields []string) (r []MExpr) {
	for _, f := range fields {
		r = append(r, ME_Field(f))
	}
	return
}

func ME_Field(f string) MExpr {
	return MExpr{field: f, operator: MYSQL_OP_FIELD}
}

func ME_Alias(f string, d string) MExpr {
	return MExpr{field: f, operator: MYSQL_OP_ALIAS, oprand1: d}
}

func ME_Expr_Alias(d1 MExpr, d2 string) MExpr {
	return MExpr{operator: MYSQL_OP_EXPR_ALIAS, oprand1: d1, oprand2: d2}
}
func ME_Assign(f string, d MExpr) MExpr {
	return MExpr{field: f, operator: MYSQL_OP_ASSIGN, oprand1: d}
}
func ME_Assign_Value(f string, d interface{}) MExpr {
	return MExpr{field: f, operator: MYSQL_OP_ASSIGN_VALUE, oprand1: d}
}

func ME_Eq(f string, d interface{}) MExpr {
	return MExpr{field: f, operator: MYSQL_OP_EQ, oprand1: d}
}

func ME_Ne(f string, d interface{}) MExpr {
	return MExpr{field: f, operator: MYSQL_OP_NE, oprand1: d}
}

func ME_Ge(f string, d interface{}) MExpr {
	return MExpr{field: f, operator: MYSQL_OP_GE, oprand1: d}
}

func ME_Gt(f string, d interface{}) MExpr {
	return MExpr{field: f, operator: MYSQL_OP_GT, oprand1: d}
}

func ME_Le(f string, d interface{}) MExpr {
	return MExpr{field: f, operator: MYSQL_OP_LE, oprand1: d}
}

func ME_Lt(f string, d interface{}) MExpr {
	return MExpr{field: f, operator: MYSQL_OP_LT, oprand1: d}
}

func ME_Between(f string, d1 interface{}, d2 interface{}) MExpr {
	return MExpr{field: f, operator: MYSQL_OP_BETWEEN, oprand1: d1, oprand2: d2}
}

func ME_IsNull(f string) MExpr {
	return MExpr{field: f, operator: MYSQL_OP_IS_NULL}
}

func ME_NotNull(f string) MExpr {
	return MExpr{field: f, operator: MYSQL_OP_NOT_NULL}
}

func ME_Like(f string, d string) MExpr {
	return MExpr{field: f, operator: MYSQL_OP_LIKE, oprand1: d}
}

func ME_NotLike(f string, d string) MExpr {
	return MExpr{field: f, operator: MYSQL_OP_NOT_LIKE, oprand1: d}
}

func ME_In(f string, d []interface{}) MExpr {
	return MExpr{field: f, operator: MYSQL_OP_IN, oprand1: d}
}

func ME_NotIn(f string, d []interface{}) MExpr {
	return MExpr{field: f, operator: MYSQL_OP_NOT_IN, oprand1: d}
}

func ME_Literal(d string) MExpr {
	return MExpr{operator: MYSQL_OP_LITERAL, oprand1: d}
}

func ME_And(d []MExpr) MExpr {
	return MExpr{operator: MYSQL_OP_AND, oprand1: d}
}

func ME_Or(d []MExpr) MExpr {
	return MExpr{operator: MYSQL_OP_OR, oprand1: d}
}

func ME_Not(d MExpr) MExpr {
	return MExpr{operator: MYSQL_OP_NOT, oprand1: d}
}

func (me MExpr) FlagsOrder() uint8 {
	return uint8((me.flags & MYSQL_FLAGS_ORDER_MASK) >> MYSQL_FLAGS_ORDER_OFFSET)
}

func (me MExpr) Asc() MExpr {
	me.flags = (me.flags & MYSQL_FLAGS_ORDER_UNMASK) | MYSQL_FLAGS_ORDER_ASC
	return me
}

func (me MExpr) Desc() MExpr {
	me.flags = (me.flags & MYSQL_FLAGS_ORDER_UNMASK) | MYSQL_FLAGS_ORDER_DESC
	return me
}
