package orm

import (
	"errors"
)

var (
	ErrUnknownType = errors.New("unknown statement type")
)

type Statement struct {
	God *God

	SelectExprs     []*Expr
	UpdateExprs     []*Expr
	InsertIntoExprs []*Expr
	DeleteExprs     []*Expr

	SetExprs     []*Expr
	ColsList     []string
	ValuesExprs  [][]*Expr
	FromExprs    []*Expr
	WhereExprs   []*Expr
	GroupbyExprs []*Expr
	HavingExprs  []*Expr
	OrderbyExprs []*Expr
	OndupExprs   []*Expr
	LimitExpr    *Expr
}

func (s *Statement) Select(model interface{}, fields ...string) error {
	if len(fields) == 0 {
		return g.Select_raw(model, E_field("*"))
	}
	field_exprs := make([]*Expr, len(fields))
	for i, f := range fields {
		field_exprs[i] = E_field(f)
	}
	return g.Select_raw(model, field_exprs...)
}

func (g *God) Select_raw(exprs ...*Expr) error {
	return (&Statement{}).Select(exprs...)
}

func Select_raw(exprs ...*Expr) *Statement {
	return (&Statement{}).Select(exprs...)
}

func Update(exprs ...*Expr) *Statement {
	return (&Statement{}).Update(exprs...)
}

func InsertInto(exprs ...*Expr) *Statement {
	return (&Statement{}).InsertInto(exprs...)
}

func Delete(exprs ...*Expr) *Statement {
	return (&Statement{}).Delete(exprs...)
}

func Set(exprs ...*Expr) *Statement {
	return (&Statement{}).Set(exprs...)
}

func Cols(cols ...string) *Statement {
	return (&Statement{}).Cols(cols...)
}

func Values(exprs ...[]*Expr) *Statement {
	return (&Statement{}).Values(exprs...)
}

func From(exprs ...*Expr) *Statement {
	return (&Statement{}).From(exprs...)
}

func Where(exprs ...*Expr) *Statement {
	return (&Statement{}).Where(exprs...)
}

func Groupby(exprs ...*Expr) *Statement {
	return (&Statement{}).Groupby(exprs...)
}

func Having(exprs ...*Expr) *Statement {
	return (&Statement{}).Having(exprs...)
}

func Orderby(exprs ...*Expr) *Statement {
	return (&Statement{}).Orderby(exprs...)
}

func Ondup(exprs ...*Expr) *Statement {
	return (&Statement{}).Ondup(exprs...)
}

func Limit(args ...int) *Statement {
	return (&Statement{}).Limit(args...)
}

func (s *Statement) Select(exprs ...*Expr) *Statement {
	s.Type = STATEMENT_TYPE_SELECT
	s.SelectExprs = append(s.SelectExprs, exprs...)
	return s
}

func (s *Statement) Update(exprs ...*Expr) *Statement {
	s.Type = STATEMENT_TYPE_UPDATE
	s.UpdateExprs = append(s.UpdateExprs, exprs...)
	return s
}

func (s *Statement) InsertInto(exprs ...*Expr) *Statement {
	s.Type = STATEMENT_TYPE_INSERT
	s.InsertIntoExprs = append(s.InsertIntoExprs, exprs...)
	return s
}

func (s *Statement) Delete(exprs ...*Expr) *Statement {
	s.Type = STATEMENT_TYPE_DELETE
	s.DeleteExprs = append(s.DeleteExprs, exprs...)
	return s
}

func (s *Statement) Set(exprs ...*Expr) *Statement {
	s.SetExprs = append(s.SetExprs, exprs...)
	return s
}

func (s *Statement) Cols(cols ...string) *Statement {
	s.ColsList = append(s.ColsList, cols...)
	return s
}

func (s *Statement) Values(exprs ...[]*Expr) *Statement {
	s.ValuesExprs = append(s.ValuesExprs, exprs...)
	return s
}

func (s *Statement) Values_append(exprs ...*Expr) *Statement {
	s.ValuesExprs = append(s.ValuesExprs, exprs)
	return s
}

func (s *Statement) From(exprs ...*Expr) *Statement {
	s.FromExprs = append(s.FromExprs, exprs...)
	return s
}

func (s *Statement) Where(exprs ...*Expr) *Statement {
	s.WhereExprs = append(s.WhereExprs, exprs...)
	return s
}

func (s *Statement) Groupby(exprs ...*Expr) *Statement {
	s.GroupbyExprs = append(s.GroupbyExprs, exprs...)
	return s
}

func (s *Statement) Having(exprs ...*Expr) *Statement {
	s.HavingExprs = append(s.HavingExprs, exprs...)
	return s
}

func (s *Statement) Orderby(exprs ...*Expr) *Statement {
	s.OrderbyExprs = append(s.OrderbyExprs, exprs...)
	return s
}

func (s *Statement) Ondup(exprs ...*Expr) *Statement {
	s.OndupExprs = append(s.OndupExprs, exprs...)
	return s
}

func (s *Statement) Limit(args ...int) *Statement {
	if len(args) == 1 {
		s.LimitExpr = E_literal("?", args[0])
	} else if len(args) >= 2 {
		s.LimitExpr = E_literal("? OFFSET ?", args[0], args[1])
	}
	return s
}
