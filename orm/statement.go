package orm

import (
	"database/sql"
	"reflect"
)

type insert_clause struct {
	cols   []string
	values []*Expr
	sets   []*Expr
}

type Statement struct {
	God *God

	SQL   string
	Args  SqlArgs
	Error error

	SelectClause []*Expr
	UpdateClause []*Expr
	DeleteClause []*Expr
	InsertClause *insert_clause

	TableClause  []*Expr
	ShardingData [][]interface{}

	OnClause      []*Expr
	GroupbyClause []*Expr
	HavingClause  []*Expr
	OrderbyClause []*Expr
	LimitClause   *Expr
	OndupClause   []*Expr
}

func (s *Statement) Count() (int, error) {
	if m, err := s.One(E_field("COUNT(*)").Alias("count")); err != nil {
		return 0, err
	} else if m != nil {
		return m.Int("count"), nil
	} else {
		return 0, nil
	}
}

/*
One(&user, "id", "name", E_field("age + 1").As("age"), E_field("status = 'del'").As("is_del"))
One(&user, "tpl_simple")
*/
func (s *Statement) One(args ...interface{}) (Model, error) {
	s.SelectClause = s.God.args_to_field_exprs_with_tpl(args...)
	if len(s.TableClause) == 0 {
		s.TableClause = []*Expr{E_table(s.God.Table.Name())}
	}
	s.Limit(1)

	if err := s.parse_select(); err != nil {
		return nil, err
	}

	rows, err := s.Query(s.SQL, s.Args...)
	if err != nil {
		return nil, err
	}

	return s.God.rows_to_model(rows)
}

/*
Multi(&users, "id", "name", E_field("age + 1").As("age"), E_field("status = refused").As("is_del"))
Multi(&users, "tpl_full")
*/
func (s *Statement) Multi(args ...interface{}) ([]Model, error) {
	s.SelectClause = s.God.args_to_field_exprs_with_tpl(args...)
	if len(s.TableClause) == 0 {
		s.TableClause = []*Expr{E_table(s.God.Table.Name())}
	}

	if err := s.parse_select(); err != nil {
		return nil, err
	}
	rows, err := s.Query(s.SQL, s.Args...)
	if err != nil {
		return nil, err
	}
	return s.God.rows_to_models(rows)
}

/*
Update(&user)
	=> UPDATE <table> SET <non_nil_field_name_1> = <non_nil_field_val_1>, <non_nil_field_name_2> = <non_nil_field_val_2> .. WHERE <pk_field_name_1> = <pk_field_val_1> AND <pk_field_name_2> = <pk_field_val_2>
Update(E_assign("name", "Goosman-lei"), E_assign("age", 31))
	=> UPDATE <table> SET name = 'Goosman-lei', age = 31
*/
func (s *Statement) Update(args ...interface{}) (int, error) {
	on_exprs, set_exprs := s.God.parse_args_for_update(args...)
	if len(on_exprs) > 0 {
		s.OnClause = on_exprs
	}
	if len(set_exprs) > 0 {
		s.UpdateClause = set_exprs
	}
	if len(s.TableClause) == 0 {
		s.TableClause = []*Expr{E_table(s.God.Table.Name())}
	}

	if err := s.parse_update(); err != nil {
		return 0, err
	} else if r, err := s.Exec(s.SQL, s.Args...); err != nil {
		return 0, err
	} else if cnt, err := r.RowsAffected(); err != nil {
		return 0, err
	} else {
		return int(cnt), nil
	}
}

/*
Delete(&user)
	=> DELETE FROM <table> WHERE <pk_field_name_1> = <pk_field_val_1> AND <pk_field_name_2> = <pk_field_val_2>
Delete(pk...)
	=> DELETE FROM <table> WHERE <pk_field_name_1> = <pk_field_val_1> AND <pk_field_name_2> = <pk_field_val_2>
Delete(E_eq("name", "Goosman-lei"), E_le("age", 100))
	=> DELETE FROM <table> WHERE name = 'Goosman-lei' AND age < 100
*/
func (s *Statement) Delete(args ...interface{}) (int, error) {
	s.DeleteClause = s.God.parse_args_for_delete(args...)
	if len(s.TableClause) == 0 {
		s.TableClause = []*Expr{E_table(s.God.Table.Name())}
	}

	if err := s.parse_delete(); err != nil {
		return 0, err
	} else if r, err := s.Exec(s.SQL, s.Args...); err != nil {
		return 0, err
	} else if cnt, err := r.RowsAffected(); err != nil {
		return 0, err
	} else {
		return int(cnt), nil
	}
}

/*
Insert(&user, &user, ...)
	=> INSERT INTO <table> (non_nil_field_name...) VALUES(non_nil_field_val...), (non_nil_field_val...)...
Insert(E_assign("name", "Goosman-lei"), E_assign("age", 31))
	=> INSERT INTO <table> SET name = 'Goosman-lei', age = 31
Insert(E_values("Goosman-lei", 31), E_values("Jacky", 28))
	=> INSERT INTO <table> VALUES("Goosman-lei", 31), ("Jacky", 28)
Insert(E_fields("name", "age"), E_values("Goosman-lei", 31), E_values("Jacky", 28))
	=> INSERT INTO <table> (name, age) VALUES("Goosman-lei", 31), ("Jacky", 28)
*/
func (s *Statement) Insert(args ...interface{}) (int, error) {
	s.InsertClause = s.God.parse_args_for_insert(args...)
	if len(s.TableClause) == 0 {
		s.TableClause = []*Expr{E_table(s.God.Table.Name())}
	}

	if err := s.parse_insert(); err != nil {
		return 0, err
	} else if r, err := s.Exec(s.SQL, s.Args...); err != nil {
		return 0, err
	} else if last_id, err := r.LastInsertId(); err != nil {
		return 0, err
	} else {
		return int(last_id), nil
	}
}

/*
On(E_eq("name", "Goosman-lei"), E_gt("age", 30))
	=> WHERE name = 'Goosman-lei' AND age > 30
*/
func (s *Statement) On(args ...*Expr) *Statement {
	s.OnClause = args
	return s
}

/*
Groupby("age", E_field("id").Desc())
	=> GROUP BY age, id DESC
*/
func (s *Statement) Groupby(args ...interface{}) *Statement {
	s.GroupbyClause = s.God.args_to_field_exprs(args...)
	return s
}

/*
On(E_eq("name", "Goosman-lei"), E_gt("age", 30))
	=> HAVING name = 'Goosman-lei' AND age > 30
*/
func (s *Statement) Having(args ...*Expr) *Statement {
	s.HavingClause = args
	return s
}

/*
Orderby("age", E_field("id").Desc())
	=> ORDERBY age, id DESC
*/
func (s *Statement) Orderby(args ...interface{}) *Statement {
	s.OrderbyClause = s.God.args_to_field_exprs(args...)
	return s
}

/*
Ondup(E_assign("name", "Goosman-lei"), E_assign("age", 31))
	=> ON DUPLICATE KEY UPDATE name = 'Goosman-lei', age = 31
*/
func (s *Statement) Ondup(args ...*Expr) *Statement {
	s.OndupClause = args
	return s
}

/*
From(E_table("kk_user").Alias("u").Join(E_table("kk_user_show").Alias("us"), E_literal("u.id = us.uid")))
	=> ON DUPLICATE KEY UPDATE name = 'Goosman-lei', age = 31
*/
func (s *Statement) From(args ...*Expr) *Statement {
	s.TableClause = args
	return s
}

/*
Limit(20)
	=> LIMIT 20
Limit(20, 5)
	=> LIMIT 20 OFFSET 5
*/
func (s *Statement) Limit(args ...int) *Statement {
	if len(args) == 1 {
		s.LimitClause = E_literal("?", args[0])
	} else if len(args) == 2 {
		s.LimitClause = E_literal("? OFFSET ?", args[0], args[1])
	}
	return s
}

func (s *Statement) Query(query string, args ...interface{}) (*sql.Rows, error) {
	s.God.LastStatement = s
	return s.God.query(query, args...)
}

func (s *Statement) Exec(query string, args ...interface{}) (sql.Result, error) {
	s.God.LastStatement = s
	return s.God.exec(query, args...)
}

/*
Delete(&user)
	=> DELETE FROM <table> WHERE <pk_field_name_1> = <pk_field_val_1> AND <pk_field_name_2> = <pk_field_val_2>
Delete(pk...)
	=> DELETE FROM <table> WHERE <pk_field_name_1> = <pk_field_val_1> AND <pk_field_name_2> = <pk_field_val_2>
Delete(E_eq("name", "Goosman-lei"), E_le("age", 100))
	=> DELETE FROM <table> WHERE name = 'Goosman-lei' AND age < 100
*/
func (g *God) parse_args_for_delete(args ...interface{}) []*Expr {
	if len(args) == 0 {
		return Exprs_empty
	}

	exprs := []*Expr{}
	if len(args) == 1 && g.is_model(args[0]) {
		if len(g.Model.PK) < 1 {
			return Exprs_empty
		}
		r_model := reflect.ValueOf(args[0]).Elem()
		for _, pk := range g.Model.PK {
			r_field := r_model.FieldByName(pk.R_StructField.Name)
			exprs = append(exprs, E_eq(pk.Column, r_field.Interface()))
		}
	} else if _, ok := args[0].(*Expr); !ok {
		if len(g.Model.PK) != len(args) {
			return Exprs_empty
		}
		for i, pk := range g.Model.PK {
			exprs = append(exprs, E_eq(pk.Column, args[i]))
		}
	} else {
		for _, expr := range args {
			if v, ok := expr.(*Expr); ok {
				exprs = append(exprs, v)
			}
		}
	}
	return exprs
}

/*
Update(&user)
	=> UPDATE <table> SET <non_nil_field_name_1> = <non_nil_field_val_1>, <non_nil_field_name_2> = <non_nil_field_val_2> .. WHERE <pk_field_name_1> = <pk_field_val_1> AND <pk_field_name_2> = <pk_field_val_2>
Update(&user, []string{"name", "desc"})
	=> UPDATE <table> SET <non_nil_field_name_1> = <non_nil_field_val_1>, <non_nil_field_name_2> = <non_nil_field_val_2> .. WHERE <pk_field_name_1> = <pk_field_val_1> AND <pk_field_name_2> = <pk_field_val_2>
Update(E_assign("name", "Goosman-lei"), E_assign("age", 31))
	=> UPDATE <table> SET name = 'Goosman-lei', age = 31
*/
func (g *God) parse_args_for_update(args ...interface{}) ([]*Expr, []*Expr) {
	if len(args) == 0 {
		return Exprs_empty, Exprs_empty
	}

	on_exprs := []*Expr{}
	set_exprs := []*Expr{}
	if g.is_model(args[0]) {
		cols_hints := map[string]bool{}
		if v, ok := args[len(args)-1].([]string); ok {
			for _, arg := range v {
				cols_hints[arg] = true
			}
		}
		r_model := reflect.ValueOf(args[0]).Elem()
		pk_fields := map[string]bool{}
		if len(g.Model.PK) > 0 {
			for _, pk := range g.Model.PK {
				r_field := r_model.FieldByName(pk.R_StructField.Name)
				pk_fields[pk.R_StructField.Name] = true
				on_exprs = append(on_exprs, E_eq(pk.Column, r_field.Interface()))
			}
		}
		// extract pk information and buid on clause
		for i := 0; i < r_model.NumField(); i++ {
			if v, ok := g.Model.Mapping_with_index[i]; ok && !pk_fields[r_model.Type().Field(i).Name] && (len(cols_hints) == 0 || cols_hints[v.Column]) {
				set_exprs = append(set_exprs, E_assign(v.Column, r_model.Field(i).Interface()))
			}
		}
		return on_exprs, set_exprs
	}

	for _, arg := range args {
		if v, ok := arg.(*Expr); ok {
			set_exprs = append(set_exprs, v)
		}
	}
	return Exprs_empty, set_exprs
}

/*
Insert(&user, &user, ...)
	=> INSERT INTO <table> (non_nil_field_name...) VALUES(non_nil_field_val...), (non_nil_field_val...)...
Insert(&user, &user, ..., []string{"name", "desc"})
	=> INSERT INTO <table> (non_nil_field_name...) VALUES(non_nil_field_val...), (non_nil_field_val...)...
Insert(E_assign("name", "Goosman-lei"), E_assign("age", 31))
	=> INSERT INTO <table> SET name = 'Goosman-lei', age = 31
Insert(E_values("Goosman-lei", 31), E_values("Jacky", 28))
	=> INSERT INTO <table> VALUES("Goosman-lei", 31), ("Jacky", 28)
Insert([]string{"name", "age"}, E_values("Goosman-lei", 31), E_values("Jacky", 28))
	=> INSERT INTO <table> (name, age) VALUES("Goosman-lei", 31), ("Jacky", 28)
*/
func (g *God) parse_args_for_insert(args ...interface{}) *insert_clause {
	if len(args) == 0 {
		return &insert_clause{}
	}

	cols := []string{}
	set_exprs := []*Expr{}
	value_exprs := []*Expr{}
	if g.is_model(args[0]) {
		cols_hints := map[string]bool{}
		if v, ok := args[len(args)-1].([]string); ok {
			args = args[:len(args)-1]
			for _, arg := range v {
				cols_hints[arg] = true
			}
		} else {
			for col, _ := range g.Model.Mapping_with_column {
				cols_hints[col] = true
			}
		}

		for idx, arg := range args {
			r_model := reflect.ValueOf(arg).Elem()
			vals := []interface{}{}
			for i := 0; i < r_model.NumField(); i++ {
				if field, ok := g.Model.Mapping_with_index[i]; ok && (len(cols_hints) == 0 || cols_hints[field.Column]) {
					if idx == 0 {
						cols = append(cols, field.Column)
					}
					vals = append(vals, r_model.FieldByName(field.R_StructField.Name).Interface())
				}
			}
			value_exprs = append(value_exprs, E_values(vals...))
		}
	} else if v, ok := args[0].(*Expr); ok && v.op == OP_ASSIGN {
		for _, e := range args {
			if v, ok := e.(*Expr); ok && v.op == OP_ASSIGN {
				set_exprs = append(set_exprs, v)
			}
		}
	} else if ok && v.op == OP_VALUES {
		for _, e := range args {
			if v, ok := e.(*Expr); ok && v.op == OP_VALUES {
				value_exprs = append(value_exprs, v)
			}
		}
	} else if v, ok := args[0].([]string); ok {
		cols = v
		for _, e := range args[1:] {
			if v, ok := e.(*Expr); ok && v.op == OP_VALUES {
				value_exprs = append(value_exprs, v)
			}
		}
	}
	return &insert_clause{
		cols:   cols,
		sets:   set_exprs,
		values: value_exprs,
	}
}
