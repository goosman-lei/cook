package orm

import (
	"database/sql"
	cook_conn "gitlab.niceprivate.com/golang/cook/connector"
)

func (g *God) query(query string, args ...interface{}) (*sql.Rows, error) {
	if db, err := cook_conn.GetMysql(g.Node); err != nil {
		return nil, nil
	} else {
		return db.Query(query, args...)
	}
}

func (g *God) exec(query string, args ...interface{}) (sql.Result, error) {
	if db, err := cook_conn.GetMysql(g.Node); err != nil {
		return nil, nil
	} else {
		return db.Exec(query, args...)
	}
}
