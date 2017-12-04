package orm

import (
	"gitlab.niceprivate.com/golang/cook/orm-ng/table"
)

type God struct {
	Model *Ref_Model
	Node  string
	Table table.Table
}

func NewGod(model interface{}, node string, table table.Table) *God {
	return &God{
		Model: NewRefModel(model),
		Node:  node,
		Table: table,
	}
}
