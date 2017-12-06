package orm

import (
	"database/sql"
	"errors"
)

var (
	Err_invalid_pk_for_load = errors.New("In Load() and Loads(), Model must have and only have one pk")
)

/*
Load(&user, 1)
*/
func (g *God) Load(pkv interface{}) (Model, error) {
	if len(g.Model.PK) != 1 {
		return nil, Err_invalid_pk_for_load
	}
	return g.NewStatement().On(E_eq(g.Model.PK[0].Column, pkv)).One()
}

/*
Loads(&users, 1, 2, 3)
*/
func (g *God) Loads(pkvs ...interface{}) ([]Model, error) {
	if len(g.Model.PK) != 1 {
		return nil, Err_invalid_pk_for_load
	}

	if len(pkvs) < 1 {
		return []Model{}, nil
	}

	return g.NewStatement().On(E_in(g.Model.PK[0].Column, pkvs)).Multi()
}

func (g *God) Count() (int, error) {
	return g.NewStatement().Count()
}

/*
One(&user, "id", "name", E_field("age + 1").As("age"), E_field("status = 'del'").As("is_del"))
	=> SELECT id, name, age + 1 AS age, status = 'del' AS is_del

Available Decorator:
On
Groupby
Having
Orderby
Limit
*/
func (g *God) One(args ...interface{}) (Model, error) {
	return g.NewStatement().One(args...)
}

/*
Multi(&users, "id", "name", E_field("age + 1").As("age"), E_field("status = refused").As("is_del"))
	=> SELECT id, name, age + 1 AS age, status = 'del' AS is_del

Available Decorator:
On
Groupby
Having
Orderby
Limit
*/
func (g *God) Multi(args ...interface{}) ([]Model, error) {
	return g.NewStatement().Multi(args...)
}

/*
Update(&user)
	=> UPDATE <table> SET <non_nil_field_name_1> = <non_nil_field_val_1>, <non_nil_field_name_2> = <non_nil_field_val_2> .. WHERE <pk_field_name_1> = <pk_field_val_1> AND <pk_field_name_2> = <pk_field_val_2>
Update(E_assign("name", "Goosman-lei"), E_assign("age", 31))
	=> UPDATE <table> SET name = 'Goosman-lei', age = 31

Available Decorator:
On
Orderby
Limit
*/
func (g *God) Update(args ...interface{}) (int, error) {
	return g.NewStatement().Update(args...)
}

/*
Insert(&user)
	=> INSERT INTO <table> (non_nil_field_name...) VALUES(non_nil_field_val...)
Insert(&users)
	=> INSERT INTO <table> (non_nil_field_name...) VALUES(non_nil_field_val...), (non_nil_field_val...)...
Insert(E_assign("name", "Goosman-lei"), E_assign("age", 31))
	=> INSERT INTO <table> SET name = 'Goosman-lei', age = 31
Insert(E_fields("name", "age"), E_values("Goosman-lei", 31), E_values("Jacky", 28))
	=> INSERT INTO <table> (name, age) VALUES("Goosman-lei", 31), ("Jacky", 28)

Available Decorator:
Ondup
*/
func (g *God) Insert(args ...interface{}) (int, error) {
	return g.NewStatement().Insert(args...)
}

/*
Delete(&user)
	=> DELETE FROM <table> WHERE <pk_field_name_1> = <pk_field_val_1> AND <pk_field_name_2> = <pk_field_val_2>
Delete(pk...)
	=> DELETE FROM <table> WHERE <pk_field_name_1> = <pk_field_val_1> AND <pk_field_name_2> = <pk_field_val_2>
Delete(E_eq("name", "Goosman-lei"), E_le("age", 100))
	=> DELETE FROM <table> WHERE name = 'Goosman-lei' AND age < 100

Available Decorator:
Orderby
Limit
*/
func (g *God) Delete(args ...interface{}) (int, error) {
	return g.NewStatement().Delete(args...)
}

/*
On(E_eq("name", "Goosman-lei"), E_gt("age", 30))
	=> WHERE name = 'Goosman-lei' AND age > 30
*/
func (g *God) On(args ...*Expr) *Statement {
	return g.NewStatement().On(args...)
}

/*
Groupby("age", E_field("id").Desc())
	=> GROUP BY age, id DESC
*/
func (g *God) Groupby(args ...interface{}) *Statement {
	return g.NewStatement().Groupby(args...)
}

/*
Having(E_eq("name", "Goosman-lei"), E_gt("age", 30))
	=> HAVING name = 'Goosman-lei' AND age > 30
*/
func (g *God) Having(args ...*Expr) *Statement {
	return g.NewStatement().Having(args...)
}

/*
Orderby("age", E_field("id").Desc())
	=> ORDERBY age, id DESC
*/
func (g *God) Orderby(args ...interface{}) *Statement {
	return g.NewStatement().Orderby(args...)
}

/*
Ondup(E_assign("name", "Goosman-lei"), E_assign("age", 31))
	=> ON DUPLICATE KEY UPDATE name = 'Goosman-lei', age = 31
*/
func (g *God) Ondup(args ...*Expr) *Statement {
	return g.NewStatement().Ondup(args...)
}

/*
Limit(20)
	=> LIMIT 20
Limit(20, 5)
	=> LIMIT 20 OFFSET 5
*/
func (g *God) Limit(args ...int) *Statement {
	return g.NewStatement().Limit(args...)
}

func (g *God) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return g.NewStatement().Query(query, args...)
}

func (g *God) Exec(query string, args ...interface{}) (sql.Result, error) {
	return g.NewStatement().Exec(query, args...)
}
