package mysql

import (
	"testing"
)

func Test_parse_select_1(t *testing.T) {
	var (
		q     *Statement
		query string
		args  sql_args
		err   error
	)

	q = Select(
		E_fields("id", "name", "age", "add_time")...,
	).Select(
		E_literal("(age + 100) / 2").Alias("some_data"),
	).From(
		E_table("kk_user").Alias("u"),
	).Where(
		E_or(
			E_lt("id", 10),
			E_gt("id", 100),
		),
	).Groupby(
		E_field("age").Desc(),
	).Having(
		E_gt("add_time", 0),
	).Orderby(
		E_field("id").Asc(),
	).Limit(10, 0)

	if query, args, err = q.Parse(); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "SELECT id, name, age, add_time, (age + 100) / 2 AS some_data FROM kk_user AS u WHERE (id < ? OR id > ?) GROUP BY age DESC HAVING add_time > ? ORDER BY id ASC LIMIT ? OFFSET ?" {
		t.Logf("wrong sql: %s", query)
		t.Fail()
	} else {
		expect_args := []interface{}{10, 100, 0, 10, 0}
		for i, v := range expect_args {
			if v != args[i] {
				t.Logf("wrong args at index[%d]. want: %#v, real: %#v", i, v, args[i])
				t.Fail()
			}
		}
	}
}

func Test_parse_select_2(t *testing.T) {
	var (
		q     *Statement
		query string
		args  sql_args
		err   error
	)

	q = Select(
		E_field("u.id").Alias("uid"),
		E_field("u.name").Alias("uname"),
		E_field("s.id").Alias("sid"),
		E_field("s.pic").Alias("pic_url"),
		E_field("s.add_time").Alias("publish_time"),
	).From(
		E_table("kk_user").Alias("u").Join(
			E_table("kk_user_show").Alias("s"),
			E_literal("u.id = s.uid"),
		),
	).Where(
		E_in("u.id", []int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
		E_in("s.status", []string{"wait", "hide"}),
	).Orderby(E_field("s.id").Desc()).Limit(10, 0)

	if query, args, err = q.Parse(); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "SELECT u.id AS uid, u.name AS uname, s.id AS sid, s.pic AS pic_url, s.add_time AS publish_time FROM kk_user AS u JOIN kk_user_show AS s ON u.id = s.uid WHERE u.id IN(?, ?, ?, ?, ?, ?, ?, ?, ?) AND s.status IN(?, ?) ORDER BY s.id DESC LIMIT ? OFFSET ?" {
		t.Logf("wrong sql: %s", query)
		t.Fail()
	} else {
		expect_args := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, "wait", "hide", 10, 0}
		for i, v := range expect_args {
			if v != args[i] {
				t.Logf("wrong args at index[%d]. want: %#v, real: %#v", i, v, args[i])
				t.Fail()
			}
		}
	}
}

func Test_parse_update(t *testing.T) {
	var (
		q     *Statement
		query string
		args  sql_args
		err   error
	)

	q = Update(E_table("kk_user")).Set(
		E_literal("limit_time = limit_time + 1"),
		E_literal("last_update_time = unix_timestamp(now())"),
	).Where(
		E_lt("id", 100),
	).Orderby(
		E_field("id").Desc(),
	).Limit(10)

	if query, args, err = q.Parse(); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "UPDATE kk_user SET limit_time = limit_time + 1, last_update_time = unix_timestamp(now()) WHERE id < ? ORDER BY id DESC LIMIT ?" {
		t.Logf("wrong sql: %s", query)
		t.Fail()
	} else {
		expect_args := []interface{}{100, 10}
		for i, v := range expect_args {
			if v != args[i] {
				t.Logf("wrong args at index[%d]. want: %#v, real: %#v", i, v, args[i])
				t.Fail()
			}
		}
	}
}

func Test_parse_delete(t *testing.T) {
	var (
		q     *Statement
		query string
		args  sql_args
		err   error
	)

	q = Delete(E_table("kk_user")).Where(
		E_lt("id", 100),
	).Orderby(
		E_field("id").Desc(),
	).Limit(10)

	if query, args, err = q.Parse(); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "DELETE FROM kk_user WHERE id < ? ORDER BY id DESC LIMIT ?" {
		t.Logf("wrong sql: %s", query)
		t.Fail()
	} else {
		expect_args := []interface{}{100, 10}
		for i, v := range expect_args {
			if v != args[i] {
				t.Logf("wrong args at index[%d]. want: %#v, real: %#v", i, v, args[i])
				t.Fail()
			}
		}
	}
}

func Test_parse_insert(t *testing.T) {
	var (
		q     *Statement
		query string
		args  sql_args
		err   error
	)

	q = Insert(E_table("kk_user")).Set(
		E_assign("name", "goosman-lei"),
		E_assign("add_time", E_literal("unix_timestamp(now())")),
		E_assign("avatar_url", "http://cdn.host.com/name/avatar.png"),
		E_assign("age", 30),
	).Ondup(
		E_assign("last_update_time", E_literal("unix_timestamp(now())")),
		E_assign("avatar_url", "http://cdn.host.com/name/avatar.png"),
		E_assign("age", 30),
	)

	if query, args, err = q.Parse(); err != nil {
		t.Logf("parse error: %s", err)
		t.Fail()
	} else if query != "INSERT INTO kk_user SET name = ?, add_time = unix_timestamp(now()), avatar_url = ?, age = ? ON DUPLICATE KEY UPDATE last_update_time = unix_timestamp(now()), avatar_url = ?, age = ?" {
		t.Logf("wrong sql: %s", query)
		t.Fail()
	} else {
		expect_args := []interface{}{"goosman-lei", "http://cdn.host.com/name/avatar.png", 30, "http://cdn.host.com/name/avatar.png", 30}
		for i, v := range expect_args {
			if v != args[i] {
				t.Logf("wrong args at index[%d]. want: %#v, real: %#v", i, v, args[i])
				t.Fail()
			}
		}
	}
}
