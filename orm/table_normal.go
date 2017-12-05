package orm

type table_normal struct {
	Nam string
}

func Table_normal(name string) Table {
	return &table_normal{Nam: name}
}

func (t *table_normal) Name(cols ...interface{}) string {
	return t.Nam
}
