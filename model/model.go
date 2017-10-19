package model

import (
	"database/sql"
	cook_sql "gitlab.niceprivate.com/golang/cook/sql/mysql"
	"reflect"
)

/*
1. 集群管理
2. 分表管理: PartitionContext
3. 对象映射
*/

func (g *God) One(s *cook_sql.Statement) (interface{}, error) {
	var (
		query     string
		queryArgs cook_sql.SqlArgs
		err       error
		rows      *sql.Rows
		cols      []string

		fName string
		argv  []interface{}
	)
	if query, queryArgs, err = s.Parse(); err != nil {
		return nil, err
	}

	if rows, err = g.DB.Query(query, queryArgs...); err != nil {
		return nil, err
	}

	if cols, err = rows.Columns(); err != nil {
		return nil, err
	}

	m := g.Factory()
	for _, col := range cols {
		fName = g.Model.DBFieldMapping[col].Name
		argv = append(argv, reflect.ValueOf(m).Elem().FieldByName(fName).Addr().Interface())
	}
	rows.Next()
	if err = rows.Scan(argv...); err != nil {
		return nil, err
	}

	return m, nil
}
