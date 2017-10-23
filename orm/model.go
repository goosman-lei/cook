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
		m     interface{}
	)
	if query, queryArgs, err = s.Parse(); err != nil {
		return nil, err
	}

	if rows, err = g.DB.Query(query, queryArgs...); err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, nil
	}

	if cols, err = rows.Columns(); err != nil {
		return nil, err
	}

	m = g.Factory()
	for _, col := range cols {
		if v, ok := g.Model.DB_field_mapping[col]; !ok {
			// ignore no mapping cols
			continue
		} else {
			fName = v.Model_field_name
		}
		argv = append(argv, reflect.ValueOf(m).Elem().FieldByName(fName).Addr().Interface())
	}

	if err = rows.Scan(argv...); err != nil {
		return nil, err
	}

	return m, nil
}

func (g *God) Multi(s *cook_sql.Statement) (interface{}, error) {
	var (
		query     string
		queryArgs cook_sql.SqlArgs
		err       error
		rows      *sql.Rows
		cols      []string

		fName   string
		argv    []interface{}
		m_slice []interface{}
		m       interface{}
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

	m_slice = make([]interface{}, 0)
	for rows.Next() {
		m = g.Factory()
		argv = nil
		for _, col := range cols {
			if v, ok := g.Model.DB_field_mapping[col]; !ok {
				// ignore no mapping cols
				continue
			} else {
				fName = v.Model_field_name
			}
			argv = append(argv, reflect.ValueOf(m).Elem().FieldByName(fName).Addr().Interface())
		}

		if err = rows.Scan(argv...); err != nil {
			return nil, err
		}
		m_slice = append(m_slice, m)
	}

	return m_slice, nil
}

func (g *God) Update(s *cook_sql.Statement) (int, error) {
	var (
		query         string
		queryArgs     cook_sql.SqlArgs
		err           error
		result        sql.Result
		affected_rows int64
	)
	if query, queryArgs, err = s.Parse(); err != nil {
		return 0, err
	}

	if result, err = g.DB.Exec(query, queryArgs...); err != nil {
		return 0, err
	}

	affected_rows, err = result.RowsAffected()
	return int(affected_rows), err
}

func (g *God) Delete(s *cook_sql.Statement) (int, error) {
	var (
		query         string
		queryArgs     cook_sql.SqlArgs
		err           error
		result        sql.Result
		affected_rows int64
	)
	if query, queryArgs, err = s.Parse(); err != nil {
		return 0, err
	}

	if result, err = g.DB.Exec(query, queryArgs...); err != nil {
		return 0, err
	}

	affected_rows, err = result.RowsAffected()
	return int(affected_rows), err
}

func (g *God) Insert(s *cook_sql.Statement) (int, error) {
	var (
		query     string
		queryArgs cook_sql.SqlArgs
		err       error
		result    sql.Result
		last_id   int64
	)
	if query, queryArgs, err = s.Parse(); err != nil {
		return 0, err
	}

	if result, err = g.DB.Exec(query, queryArgs...); err != nil {
		return 0, err
	}

	last_id, err = result.LastInsertId()
	return int(last_id), err
}
