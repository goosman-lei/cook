package orm

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
func (g *God) One(model interface{}, fields ...interface{}) error {
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
func (g *God) Multi(model []interface{}, fields ...interface{}) error {
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
func (g *God) Update(args ...interface{}) error {
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
func (g *God) Insert(args ...interface{}) error {
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
func (g *God) Delete(args ...interface{}) error {
}

/*
On(E_eq("name", "Goosman-lei"), E_gt("age", 30))
	=> WHERE name = 'Goosman-lei' AND age > 30
*/
func (g *God) On(expr ...*Expr) *Statement {
}

/*
Groupby("age", E_field("id").Desc())
	=> GROUP BY age, id DESC
*/
func (g *God) Groupby(expr ...interface{}) *Statement {
}

/*
On(E_eq("name", "Goosman-lei"), E_gt("age", 30))
	=> HAVING name = 'Goosman-lei' AND age > 30
*/
func (g *God) Having(expr ...*Expr) *Statement {
}

/*
Orderby("age", E_field("id").Desc())
	=> ORDERBY age, id DESC
*/
func (g *God) Orderby(expr ...interface{}) *Statement {
}

/*
Ondup(E_assign("name", "Goosman-lei"), E_assign("age", 31))
	=> ON DUPLICATE KEY UPDATE name = 'Goosman-lei', age = 31
*/
func (g *God) Ondup(expr ...*Expr) *Statement {
}

/*
Limit(20)
	=> LIMIT 20
Limit(20, 5)
	=> LIMIT 20 OFFSET 5
*/
func (g *God) Limit(expr ...int) *Statement {
}
