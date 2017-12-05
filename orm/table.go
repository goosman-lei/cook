package orm

import (
	"reflect"
)

type Table interface {
	Name(cols ...interface{}) string
}

type ShardingTable interface {
	Table
	Names(cols ...interface{}) map[string][][]interface{}
}

func Names(t Table, cols ...interface{}) map[string][][]interface{} {
	mapping := make(map[string][][]interface{})
	col_cnt := len(cols)
	max_cnt := 0
	cols_r := []reflect.Value{}
	for _, col_vals := range cols {
		col_vals_r := reflect.ValueOf(col_vals)
		if col_vals_r.Kind() != reflect.Slice {
			// wrap into slice
			col_vals_r = reflect.ValueOf([]interface{}{col_vals})
		}
		if col_vals_r.Len() > max_cnt {
			max_cnt = col_vals_r.Len()
		}
		cols_r = append(cols_r, col_vals_r)
	}
	for i := 0; i < max_cnt; i++ {
		args := []interface{}{}
		for j := 0; j < col_cnt; j++ {
			if cols_r[j].Len() <= 0 {
				args = append(args, nil)
			} else if cols_r[j].Len() <= i {
				args = append(args, cols_r[j].Index(cols_r[j].Len()-1).Interface())
			} else {
				args = append(args, cols_r[j].Index(i).Interface())
			}
		}
		name := t.Name(args...)
		if _, ok := mapping[name]; !ok {
			mapping[name] = make([][]interface{}, col_cnt)
			for k := 0; k < col_cnt; k++ {
				mapping[name][k] = []interface{}{args[k]}
			}
		} else {
			for k := 0; k < col_cnt; k++ {
				mapping[name][k] = append(mapping[name][k], args[k])
			}
		}
	}
	return mapping
}
