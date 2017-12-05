package orm

import (
	cook_util "gitlab.niceprivate.com/golang/cook/util"
	"reflect"
)

type God struct {
	Model *Ref_Model
	Node  string
	Table Table
	Tpls  map[string][]*Expr

	Silent bool // in silent mode, SQL will not execute. used for debug
}

var (
	Expr_star   = E_field("*")
	Exprs_star  = []*Expr{E_field("*")}
	Exprs_empty = []*Expr{}
)

func NewGod(model interface{}, node string, table Table) *God {
	return &God{
		Model:  NewRefModel(model),
		Node:   node,
		Table:  table,
		Tpls:   make(map[string][]*Expr),
		Silent: false,
	}
}

func (g *God) Tpl(tpl string, fields ...interface{}) {
	g.Tpls[tpl] = g.args_to_field_exprs(fields...)
}

func (g *God) NewStatement() *Statement {
	return &Statement{
		God: g,
	}
}

func (g *God) args_to_field_exprs(args ...interface{}) []*Expr {
	if len(args) == 0 {
		return Exprs_star
	}
	exprs := []*Expr{}
	for _, field := range args {
		switch v := field.(type) {
		case string:
			exprs = append(exprs, E_field(v))
		case *Expr:
			exprs = append(exprs, v)
		default:
			cook_util.Panicf("field type is not string and *Expr")
		}
	}
	return exprs
}

func (g *God) is_model(model interface{}) bool {
	return reflect.TypeOf(model).Kind() == reflect.Ptr && reflect.TypeOf(model).Elem() == g.Model.R_Type
}

func (g *God) NewModel() Model {
	return reflect.New(g.Model.R_Type).Interface().(Model)
}
