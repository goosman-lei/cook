package orm

import (
	"errors"
	cook_sql "gitlab.niceprivate.com/golang/cook/sql/mysql"
	"reflect"
)

var (
	ErrInvalidArgs   = errors.New("every args of Tables() must be slice")
	ErrInvalidArgLen = errors.New("every args of Tables() must be same length")
)

func (g *God) Table(input_args ...interface{}) *cook_sql.Expr {
	return cook_sql.E_table(g.shard(input_args...))
}

func (g *God) Tables(input_args ...interface{}) ([]*cook_sql.Expr, [][][]interface{}, error) {
	var (
		r_args        []reflect.Value
		flat_args     [][]interface{}
		sharding_map  map[string][][]interface{}
		sharding_name string

		final_tables []*cook_sql.Expr
		final_args   [][][]interface{}
		final_idx    int
	)

	/*
		input_args = [][]interface{}{
			[]interface{}{1, 2, 3, 4, 5}, // id
			[]interface{}{"Jack", "Tom", "Green", "John", "Linus"}, // name
		}
	*/

	// step1: check input and flat input args
	/*
		flat_args = [][]interface{
			[]interface{}{1, "Jack"},
			[]interface{}{2, "Tom"},
			[]interface{}{3, "Green"},
			[]interface{}{4, "John"},
			[]interface{}{5, "Linus"},
		}
	*/
	r_args = make([]reflect.Value, len(input_args))
	for i, arg := range input_args {
		r_args[i] = reflect.ValueOf(arg)
		if r_args[i].Kind() != reflect.Slice {
			return nil, nil, ErrInvalidArgs
		}
		if r_args[i].Len() != r_args[0].Len() {
			return nil, nil, ErrInvalidArgLen
		}
		if i == 0 {
			flat_args = make([][]interface{}, r_args[i].Len())
		}
		for j := 0; j < r_args[i].Len(); j++ {
			if i == 0 {
				flat_args[j] = []interface{}{r_args[i].Index(j).Interface()}
			} else {
				flat_args[j] = append(flat_args[j], r_args[i].Index(j).Interface())
			}
		}
	}

	// step2: apply to sharding and fill into mapping
	/*
		sharding_map = map[string][][]interface{}{
			"table_1": [][]interface{}{
				[]interface{}{1, "Jack"},
				[]interface{}{3, "Green"},
				[]interface{}{5, "Linus"},
			}
			"table_2": [][]interface{}{
				[]interface{}{2, "Tom"},
				[]interface{}{4, "John"},
			}
		}
	*/
	sharding_map = make(map[string][][]interface{})
	for _, cols := range flat_args {
		sharding_name = g.shard(cols...)
		if _, exists := sharding_map[sharding_name]; !exists {
			sharding_map[sharding_name] = [][]interface{}{cols}
		} else {
			sharding_map[sharding_name] = append(sharding_map[sharding_name], cols)
		}
	}

	// step3: restore mapping to user format
	/*
		final_tables = []*cook_sql.Expr{
			cook_sql.E_table("table_1"),
			cook_sql.E_table("table_2"),
		}
		final_args = [][][]interface{}{
			[][]interface{}{
				[]interface{}{1, 3, 5},
				[]interface{}{"Jack", "Green", "Linus"},
			},
			[][]interface{}{
				[]interface{}{2, 4},
				[]interface{}{"Tom", "John"},
			},
		}
	*/
	final_tables = make([]*cook_sql.Expr, len(sharding_map))
	final_args = make([][][]interface{}, len(sharding_map))
	final_idx = 0
	for sharding_name, args := range sharding_map {
		final_tables[final_idx] = cook_sql.E_table(sharding_name)
		for i, arg := range args {
			if i == 0 {
				final_args[final_idx] = make([][]interface{}, len(arg))
			}
			for j, col := range arg {
				if i == 0 {
					final_args[final_idx][j] = []interface{}{col}
				} else {
					final_args[final_idx][j] = append(final_args[final_idx][j], col)
				}
			}
		}
		final_idx++
	}
	return final_tables, final_args, nil
}
