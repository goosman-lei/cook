package orm

import (
	cook_conn "gitlab.niceprivate.com/golang/cook/connector"
	cook_sql "gitlab.niceprivate.com/golang/cook/sql/mysql"
	"reflect"
	"testing"
	"time"
)

type M_User struct {
	Id          int
	Name        string
	Password    string `orm:"col(passwd)"`
	Avatar      string
	Description string
}

func F_User() interface{} {
	return &M_User{}
}

func init() {
	configs := map[string]cook_conn.MysqlConf{
		"default": cook_conn.MysqlConf{
			Addr:     "10.10.200.12:3306",
			Username: "nice",
			Password: "Cb84eZaa229ddnm",
			Database: "kkgoo",

			MaxIdle:     4,
			MaxOpen:     4,
			MaxLifeTime: time.Minute * 10,
		},
	}

	cook_conn.SetupMysql(configs)
}

func Test_One(t *testing.T) {
	var (
		god  *God
		err  error
		user interface{}
	)
	if god, err = NewGod_shard_none(F_User, "default", "kk_user"); err != nil {
		t.Logf("create god failed: %s", err)
		t.Fail()
	}
	if user, err = god.One(cook_sql.Select(
		cook_sql.E_fields("id", "name", "avatar", "passwd", "description")...,
	).From(
		god.Table(),
	).Where(
		cook_sql.E_eq("id", 100),
	)); err != nil {
		t.Logf("query faield: %s", err)
		t.Fail()
	}

	if r_user := reflect.ValueOf(user); r_user.Kind() != reflect.Ptr || r_user.Elem().Kind() != reflect.Struct || r_user.Elem().Type().Name() != "M_User" || user.(*M_User).Id != 100 {
		t.Logf("want *M_User{Id: 100}, but: %#v", user)
		t.Fail()
	}
}

func Test_Multi(t *testing.T) {
	var (
		god   *God
		err   error
		users interface{}
	)
	if god, err = NewGod_shard_none(F_User, "default", "kk_user"); err != nil {
		t.Logf("create god failed: %s", err)
		t.Fail()
	}
	if users, err = god.Multi(cook_sql.Select(
		cook_sql.E_fields("id", "name", "avatar", "passwd", "description")...,
	).From(
		god.Table(),
	).Where(
		cook_sql.E_gt("id", 100),
	).Limit(10)); err != nil {
		t.Logf("query faield: %s", err)
		t.Fail()
	}

	for _, user := range users.([]interface{}) {
		if r_user := reflect.ValueOf(user); r_user.Kind() != reflect.Ptr || r_user.Elem().Kind() != reflect.Struct || r_user.Elem().Type().Name() != "M_User" || user.(*M_User).Id <= 100 {
			t.Logf("want *M_User && .Id > 100, but: %#v", user)
			t.Fail()
		}
	}
}

func Test_Insert(t *testing.T) {
	t.Skip()
	var (
		god     *God
		err     error
		last_id int
	)
	if god, err = NewGod_shard_none(F_User, "default", "kk_user"); err != nil {
		t.Logf("create god failed: %s", err)
		t.Fail()
	}

	if last_id, err = god.Insert(cook_sql.InsertInto(
		god.Table(),
	).Set(
		cook_sql.E_assign("name", "golang-model-ut"),
		cook_sql.E_assign("ctime", cook_sql.E_literal("unix_timestamp(now())")),
	)); err != nil {
		t.Logf("query faield: %s", err)
		t.Fail()
	} else if last_id <= 0 {
		t.Logf("lastid less or equal zero")
		t.Fail()
	}
}

func Test_(t *testing.T) {
	/*
		user := God_user.One(Field(E_field("*")).Where(E_eq("id", 1)))

		users := God_user.Multi(Field(

		God_user.Update(Set(
			E_assign("name", user.Name),
			E_assign("passwd", user.Passwd),
		).Where(E_eq("id", 1)))

		God_user.Delete(Where(E_eq("id", 1)))

		God_user.Insert(user.E_set("passwd"))
	*/
}
