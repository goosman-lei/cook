package orm

import (
	"testing"
)

func parse_expr(expr *Expr) (string, SqlArgs, error) {
	var (
		err   error
		args  SqlArgs  = []interface{}{}
		pargs *SqlArgs = &args

		buf []byte = make([]byte, 1024*1024)
		off int    = 0
	)

	if err = expr_to_string(expr, pargs, &buf, &off); err != nil {
		return "", nil, err
	} else {
		return string(buf[:off]), args, nil
	}
}

func Test_parse_expr_eq(t *testing.T) {
	if query, args, err := parse_expr(
		E_and(
			E_eq("age", 1),
			E_eq("age", E_field("birthday")),
		),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "(age = ? AND age = (birthday))" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	} else {
		expect_args := []interface{}{1}
		for i, v := range expect_args {
			if v != args[i] {
				t.Logf("wrong args at index[%d]. want: %#v, real: %#v", i, v, args[i])
				t.Fail()
			}
		}
	}
}

func Test_parse_expr_ne(t *testing.T) {
	if query, args, err := parse_expr(
		E_ne("age", 1),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "age != ?" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	} else {
		expect_args := []interface{}{1}
		for i, v := range expect_args {
			if v != args[i] {
				t.Logf("wrong args at index[%d]. want: %#v, real: %#v", i, v, args[i])
				t.Fail()
			}
		}
	}
}

func Test_parse_expr_ge(t *testing.T) {
	if query, args, err := parse_expr(
		E_ge("age", 1),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "age >= ?" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	} else {
		expect_args := []interface{}{1}
		for i, v := range expect_args {
			if v != args[i] {
				t.Logf("wrong args at index[%d]. want: %#v, real: %#v", i, v, args[i])
				t.Fail()
			}
		}
	}
}

func Test_parse_expr_gt(t *testing.T) {
	if query, args, err := parse_expr(
		E_gt("age", 1),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "age > ?" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	} else {
		expect_args := []interface{}{1}
		for i, v := range expect_args {
			if v != args[i] {
				t.Logf("wrong args at index[%d]. want: %#v, real: %#v", i, v, args[i])
				t.Fail()
			}
		}
	}
}

func Test_parse_expr_le(t *testing.T) {
	if query, args, err := parse_expr(
		E_le("age", 1),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "age <= ?" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	} else {
		expect_args := []interface{}{1}
		for i, v := range expect_args {
			if v != args[i] {
				t.Logf("wrong args at index[%d]. want: %#v, real: %#v", i, v, args[i])
				t.Fail()
			}
		}
	}
}

func Test_parse_expr_lt(t *testing.T) {
	if query, args, err := parse_expr(
		E_lt("age", 1),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "age < ?" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	} else {
		expect_args := []interface{}{1}
		for i, v := range expect_args {
			if v != args[i] {
				t.Logf("wrong args at index[%d]. want: %#v, real: %#v", i, v, args[i])
				t.Fail()
			}
		}
	}
}

func Test_parse_expr_between(t *testing.T) {
	if query, args, err := parse_expr(
		E_between("age", 1, 3),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "age BETWEEN ? AND ?" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	} else {
		expect_args := []interface{}{1, 3}
		for i, v := range expect_args {
			if v != args[i] {
				t.Logf("wrong args at index[%d]. want: %#v, real: %#v", i, v, args[i])
				t.Fail()
			}
		}
	}
}

func Test_parse_expr_is_null(t *testing.T) {
	if query, _, err := parse_expr(
		E_is_null("age"),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "age IS NULL" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	}
}

func Test_parse_expr_not_null(t *testing.T) {
	if query, _, err := parse_expr(
		E_not_null("age"),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "age IS NOT NULL" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	}
}

func Test_parse_expr_like(t *testing.T) {
	if query, args, err := parse_expr(
		E_like("name", "%goosman-lei"),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "name LIKE ?" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	} else {
		if v, ok := args[0].(string); !ok || v != "%goosman-lei" {
			t.Logf("wrong args: %#v", args)
			t.Fail()
		}
	}
}

func Test_parse_expr_not_like(t *testing.T) {
	if query, args, err := parse_expr(
		E_not_like("name", "%goosman-lei"),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "name NOT LIKE ?" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	} else {
		expect_args := []interface{}{"%goosman-lei"}
		for i, v := range expect_args {
			if v != args[i] {
				t.Logf("wrong args at index[%d]. want: %#v, real: %#v", i, v, args[i])
				t.Fail()
			}
		}
	}
}

func Test_parse_expr_in_slice(t *testing.T) {
	if query, args, err := parse_expr(
		E_in("age", []int{1, 2, 3, 4}),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "age IN(?, ?, ?, ?)" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	} else {
		expect_args := []interface{}{1, 2, 3, 4}
		for i, v := range expect_args {
			if v != args[i] {
				t.Logf("wrong args at index[%d]. want: %#v, real: %#v", i, v, args[i])
				t.Fail()
			}
		}
	}
}

func Test_parse_expr_in_array(t *testing.T) {
	var arr [4]string = [4]string{"goosman-lei", "goosman", "goosman.lei", "goosmanlei"}
	if query, args, err := parse_expr(
		E_in("name", arr),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "name IN(?, ?, ?, ?)" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	} else {
		expect_args := []interface{}{"goosman-lei", "goosman", "goosman.lei", "goosmanlei"}
		for i, v := range expect_args {
			if v != args[i] {
				t.Logf("wrong args at index[%d]. want: %#v, real: %#v", i, v, args[i])
				t.Fail()
			}
		}
	}
}

func Test_parse_expr_comma(t *testing.T) {
	if query, args, err := parse_expr(
		E_comma(
			E_gt("age", 10),
			E_ne("age", 20),
			E_ne("name", "goosman-lei"),
			E_lt("age", 40),
		),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "(age > ?, age != ?, name != ?, age < ?)" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	} else {
		expect_args := []interface{}{10, 20, "goosman-lei", 40}
		for i, v := range expect_args {
			if v != args[i] {
				t.Logf("wrong args at index[%d]. want: %#v, real: %#v", i, v, args[i])
				t.Fail()
			}
		}
	}
}

func Test_parse_expr_and(t *testing.T) {
	if query, args, err := parse_expr(
		E_and(
			E_gt("age", 10),
			E_ne("age", 20),
			E_ne("name", "goosman-lei"),
			E_lt("age", 40),
		),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "(age > ? AND age != ? AND name != ? AND age < ?)" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	} else {
		expect_args := []interface{}{10, 20, "goosman-lei", 40}
		for i, v := range expect_args {
			if v != args[i] {
				t.Logf("wrong args at index[%d]. want: %#v, real: %#v", i, v, args[i])
				t.Fail()
			}
		}
	}
}

func Test_parse_expr_or(t *testing.T) {
	if query, args, err := parse_expr(
		E_or(
			E_gt("age", 30),
			E_eq("age", 20),
			E_eq("name", "goosman-lei"),
			E_lt("age", 10),
		),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "(age > ? OR age = ? OR name = ? OR age < ?)" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	} else {
		expect_args := []interface{}{30, 20, "goosman-lei", 10}
		for i, v := range expect_args {
			if v != args[i] {
				t.Logf("wrong args at index[%d]. want: %#v, real: %#v", i, v, args[i])
				t.Fail()
			}
		}
	}
}

func Test_parse_expr_not(t *testing.T) {
	if query, args, err := parse_expr(
		E_not(E_gt("age", 30)),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "NOT (age > ?)" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	} else {
		expect_args := []interface{}{30}
		for i, v := range expect_args {
			if v != args[i] {
				t.Logf("wrong args at index[%d]. want: %#v, real: %#v", i, v, args[i])
				t.Fail()
			}
		}
	}
}

func Test_parse_expr_assign_value(t *testing.T) {
	if query, args, err := parse_expr(
		E_assign("name", "goosman-lei"),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "name = ?" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	} else {
		expect_args := []interface{}{"goosman-lei"}
		for i, v := range expect_args {
			if v != args[i] {
				t.Logf("wrong args at index[%d]. want: %#v, real: %#v", i, v, args[i])
				t.Fail()
			}
		}
	}
}

func Test_parse_expr_assign_expr(t *testing.T) {
	if query, args, err := parse_expr(
		E_assign("name", E_literal("concat(name, ?, ?, ?)", "goosman", "-", "lei")),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "name = (concat(name, ?, ?, ?))" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	} else {
		expect_args := []interface{}{"goosman", "-", "lei"}
		for i, v := range expect_args {
			if v != args[i] {
				t.Logf("wrong args at index[%d]. want: %#v, real: %#v", i, v, args[i])
				t.Fail()
			}
		}
	}
}

func Test_parse_expr_values_direct(t *testing.T) {
	if query, args, err := parse_expr(
		E_values(1, "goosman-lei", "password"),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "(?, ?, ?)" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	} else {
		expect_args := []interface{}{1, "goosman-lei", "password"}
		for i, v := range expect_args {
			if v != args[i] {
				t.Logf("wrong args at index[%d]. want: %#v, real: %#v", i, v, args[i])
				t.Fail()
			}
		}
	}
}

func Test_parse_expr_values_expr(t *testing.T) {
	if query, args, err := parse_expr(
		E_values(1, "goosman-lei", E_literal("unix_timestamp(now())")),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "(?, ?, (unix_timestamp(now())))" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	} else {
		expect_args := []interface{}{1, "goosman-lei"}
		for i, v := range expect_args {
			if v != args[i] {
				t.Logf("wrong args at index[%d]. want: %#v, real: %#v", i, v, args[i])
				t.Fail()
			}
		}
	}
}

func Test_parse_expr_field(t *testing.T) {
	if query, _, err := parse_expr(
		E_field("name"),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "name" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	}
}

func Test_parse_expr_literal(t *testing.T) {
	if query, args, err := parse_expr(
		E_literal("sum(age) / count(*) * ? / ?", 100, 2),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "sum(age) / count(*) * ? / ?" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	} else {
		expect_args := []interface{}{100, 2}
		for i, v := range expect_args {
			if v != args[i] {
				t.Logf("wrong args at index[%d]. want: %#v, real: %#v", i, v, args[i])
				t.Fail()
			}
		}
	}
}

func Test_parse_expr_table_normal(t *testing.T) {
	if query, _, err := parse_expr(
		E_table("kk_user").Alias("u"),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "kk_user AS u" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	}
}

func Test_parse_expr_table_join(t *testing.T) {
	if query, _, err := parse_expr(
		E_table("kk_user").Alias("u").Join(
			E_table("kk_user_show").Alias("us"),
			E_literal("u.id = us.uid"),
		),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "kk_user AS u JOIN kk_user_show AS us ON u.id = us.uid" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	}
}

func Test_parse_expr_table_left_join(t *testing.T) {
	if query, _, err := parse_expr(
		E_table("kk_user").Alias("u").Join(
			E_table("kk_user_show").Alias("us"),
			E_literal("u.id = us.uid"),
		).LeftJoin(
			E_table("kk_user_show_comment").Alias("usc"),
			E_literal("us.id = usc.sid"),
		),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "kk_user AS u JOIN kk_user_show AS us ON u.id = us.uid LEFT JOIN kk_user_show_comment AS usc ON us.id = usc.sid" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	}
}

func Test_parse_expr_table_right_join(t *testing.T) {
	if query, _, err := parse_expr(
		E_table("kk_user").Alias("u").Join(
			E_table("kk_user_show").Alias("us"),
			E_literal("u.id = us.uid"),
		).RightJoin(
			E_table("kk_user_show_comment").Alias("usc"),
			E_literal("us.id = usc.sid"),
		),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "kk_user AS u JOIN kk_user_show AS us ON u.id = us.uid RIGHT JOIN kk_user_show_comment AS usc ON us.id = usc.sid" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	}
}

func Test_parse_expr_table_inner_join(t *testing.T) {
	if query, _, err := parse_expr(
		E_table("kk_user").Alias("u").Join(
			E_table("kk_user_show").Alias("us"),
			E_literal("u.id = us.uid"),
		).InnerJoin(
			E_table("kk_user_show_comment").Alias("usc"),
			E_literal("us.id = usc.sid"),
		),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "kk_user AS u JOIN kk_user_show AS us ON u.id = us.uid INNER JOIN kk_user_show_comment AS usc ON us.id = usc.sid" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	}
}

func Test_parse_expr_table_cross_join(t *testing.T) {
	if query, _, err := parse_expr(
		E_table("kk_user").Alias("u").Join(
			E_table("kk_user_show").Alias("us"),
			E_literal("u.id = us.uid"),
		).CrossJoin(
			E_table("kk_user_show_comment").Alias("usc"),
			E_literal("us.id = usc.sid"),
		),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "kk_user AS u JOIN kk_user_show AS us ON u.id = us.uid CROSS JOIN kk_user_show_comment AS usc ON us.id = usc.sid" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	}
}

func Test_parse_expr_alias(t *testing.T) {
	if query, args, err := parse_expr(
		E_literal("sum(age) / count(*) * ? / ?", 100, 2).Alias("some_data"),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "sum(age) / count(*) * ? / ? AS some_data" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	} else {
		expect_args := []interface{}{100, 2}
		for i, v := range expect_args {
			if v != args[i] {
				t.Logf("wrong args at index[%d]. want: %#v, real: %#v", i, v, args[i])
				t.Fail()
			}
		}
	}
}

func Test_parse_expr_order(t *testing.T) {
	if query, _, err := parse_expr(
		E_comma(
			E_field("age").Asc(),
			E_field("birthday").Desc(),
		),
	); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "(age ASC, birthday DESC)" {
		t.Logf("wrong query: %s", query)
		t.Fail()
	}
}
