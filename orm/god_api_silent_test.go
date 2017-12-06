package orm

import (
	"testing"
)

type M_User_GodApiTest struct {
	*M       `orm:"nomapping"`
	Id       int `orm:"pk"`
	Name     string
	Password string `orm:"col(passwd)"`
	Avatar   string
	Desc     string `orm:"col(description)"`
	AddTime  int
	Age      int
	Sex      string
}

func F_User_GodApiTest() Model {
	return &M_User_GodApiTest{}
}

var (
	GodOf_User_ApiTest *God = NewGod(F_User_GodApiTest, "not-exists-node", Table_normal("kk_user"))
)

func init() {
	GodOf_User_ApiTest.Tpl("simple", "id", "name")
}

func Test_God_Load(t *testing.T) {
	GodOf_User_ApiTest.Load(1)
	want_sql_and_args(t, "SELECT * FROM kk_user WHERE id = ? LIMIT ?", 1, 1)

	GodOf_User_ApiTest.Loads(1, 2, 3)
	want_sql_and_args(t, "SELECT * FROM kk_user WHERE id IN(?, ?, ?)", 1, 2, 3)
}

func Test_God_Count(t *testing.T) {
	GodOf_User_ApiTest.On(
		E_between("id", 1000, 2000),
		E_gt("age", 30),
	).Count()
	want_sql_and_args(t, "SELECT COUNT(*) AS count FROM kk_user WHERE id BETWEEN ? AND ? AND age > ? LIMIT ?", 1000, 2000, 30, 1)
}

func Test_God_One(t *testing.T) {
	GodOf_User_ApiTest.One("simple")
	want_sql_and_args(t, "SELECT id, name FROM kk_user LIMIT ?", 1)

	GodOf_User_ApiTest.One("id", "name", E_literal("unix_timestamp() - add_time AS reg_secs"))
	want_sql_and_args(t, "SELECT id, name, unix_timestamp() - add_time AS reg_secs FROM kk_user LIMIT ?", 1)

	GodOf_User_ApiTest.On(E_eq("name", "Goosman-lei")).One("id", "name", E_literal("unix_timestamp() - add_time AS reg_secs"))
	want_sql_and_args(t, "SELECT id, name, unix_timestamp() - add_time AS reg_secs FROM kk_user WHERE name = ? LIMIT ?", "Goosman-lei", 1)

	GodOf_User_ApiTest.On(E_eq("sex", "male")).Groupby("age").Having(E_gt("add_time", 0)).Orderby(E_desc("id")).One("id", "name", E_literal("unix_timestamp() - add_time AS reg_secs"))
	want_sql_and_args(t, "SELECT id, name, unix_timestamp() - add_time AS reg_secs FROM kk_user WHERE sex = ? GROUP BY age HAVING add_time > ? ORDER BY id DESC LIMIT ?", "male", 0, 1)
}

func Test_God_Multi(t *testing.T) {
	GodOf_User_ApiTest.On(E_eq("sex", "male")).Groupby("age").Having(E_gt("add_time", 0)).Orderby(E_desc("id")).Limit(10, 20).Multi("id", "name", E_literal("unix_timestamp() - add_time AS reg_secs"))
	want_sql_and_args(t, "SELECT id, name, unix_timestamp() - add_time AS reg_secs FROM kk_user WHERE sex = ? GROUP BY age HAVING add_time > ? ORDER BY id DESC LIMIT ? OFFSET ?", "male", 0, 10, 20)
}

func Test_God_Update(t *testing.T) {
	user := M_User_GodApiTest{
		Id:       1,
		Name:     "Goosman-lei",
		Password: "HAHAHAHA",
		AddTime:  191919191,
		Sex:      "male",
	}

	GodOf_User_ApiTest.Update(&user)
	want_sql_and_args(t, "UPDATE kk_user SET name = ?, passwd = ?, avatar = ?, description = ?, add_time = ?, age = ?, sex = ? WHERE id = ?", "Goosman-lei", "HAHAHAHA", "", "", 191919191, 0, "male", 1)

	// hint which column will be update
	GodOf_User_ApiTest.Update(&user, []string{"name", "passwd", "add_time", "sex"})
	want_sql_and_args(t, "UPDATE kk_user SET name = ?, passwd = ?, add_time = ?, sex = ? WHERE id = ?", "Goosman-lei", "HAHAHAHA", 191919191, "male", 1)

	GodOf_User_ApiTest.On(E_lt("age", 40)).Limit(10).Update(E_literal("age = age + 1"))
	want_sql_and_args(t, "UPDATE kk_user SET age = age + 1 WHERE age < ? LIMIT ?", 40, 10)
}

func Test_God_Insert(t *testing.T) {
	users := []interface{}{
		&M_User_GodApiTest{
			Id:       1,
			Name:     "Goosman-lei",
			Password: "HAHAHAHA",
			AddTime:  191919191,
			Sex:      "male",
		},
		&M_User_GodApiTest{
			Id:       2,
			Name:     "goosman",
			Password: "HEHEHEHE",
			AddTime:  29292929,
			Sex:      "male",
		},
	}

	GodOf_User_ApiTest.Insert(users...)
	want_sql_and_args(
		t, "INSERT INTO kk_user(id, name, passwd, avatar, description, add_time, age, sex) VALUES (?, ?, ?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?, ?, ?)",
		1, "Goosman-lei", "HAHAHAHA", "", "", 191919191, 0, "male",
		2, "goosman", "HEHEHEHE", "", "", 29292929, 0, "male",
	)

	users = append(users, []string{"id", "name", "passwd", "add_time", "sex"})
	GodOf_User_ApiTest.Insert(users...)
	want_sql_and_args(
		t, "INSERT INTO kk_user(id, name, passwd, add_time, sex) VALUES (?, ?, ?, ?, ?), (?, ?, ?, ?, ?)",
		1, "Goosman-lei", "HAHAHAHA", 191919191, "male",
		2, "goosman", "HEHEHEHE", 29292929, "male",
	)

	GodOf_User_ApiTest.Ondup(
		E_assign("age", 19),
	).Insert(
		E_assign("name", "goosman"),
		E_assign("age", 18),
		E_assign("passwd", "HAHAHA"),
	)
	want_sql_and_args(t, "INSERT INTO kk_user SET name = ?, age = ?, passwd = ? ON DUPLICATE KEY UPDATE age = ?", "goosman", 18, "HAHAHA", 19)

	GodOf_User_ApiTest.Insert(E_values("goosman-lei", 18), E_values("Jacky", 28))
	want_sql_and_args(t, "INSERT INTO kk_user VALUES (?, ?), (?, ?)", "goosman-lei", 18, "Jacky", 28)

	GodOf_User_ApiTest.Insert([]string{"name", "age"}, E_values("goosman-lei", 18), E_values("Jacky", 28))
	want_sql_and_args(t, "INSERT INTO kk_user(name, age) VALUES (?, ?), (?, ?)", "goosman-lei", 18, "Jacky", 28)
}

func Test_God_Delete(t *testing.T) {
	user := M_User_GodApiTest{Id: 1}

	GodOf_User_ApiTest.Delete(&user)
	want_sql_and_args(t, "DELETE FROM kk_user WHERE id = ?", 1)

	GodOf_User_ApiTest.Delete(3)
	want_sql_and_args(t, "DELETE FROM kk_user WHERE id = ?", 3)

	GodOf_User_ApiTest.Delete(E_gt("age", 30))
	want_sql_and_args(t, "DELETE FROM kk_user WHERE age > ?", 30)
}

func want_sql_and_args(t *testing.T, sql string, args ...interface{}) {
	if GodOf_User_ApiTest.LastStatement.Error != nil {
		t.Logf("unexcept error: %s", GodOf_User_ApiTest.LastStatement.Error)
		t.Fail()
	}
	if GodOf_User_ApiTest.LastStatement.SQL != sql {
		t.Logf("unexcept sql:\n\treal:%s\n\twant:%s", GodOf_User_ApiTest.LastStatement.SQL, sql)
		t.Fail()
	}
	if len(GodOf_User_ApiTest.LastStatement.Args) != len(args) {
		t.Logf("unexcept args count:\n\treal:%d\n\twant:%d", len(GodOf_User_ApiTest.LastStatement.Args), len(args))
		t.Fail()
	}
	for i, arg := range GodOf_User_ApiTest.LastStatement.Args {
		if arg != args[i] {
			t.Logf("unexcept arg at %d:\n\treal:%#v\n\twant:%#v", i, arg, args[i])
			t.Fail()
		}
	}
}
