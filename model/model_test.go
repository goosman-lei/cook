package model

import (
	cook_conn "gitlab.niceprivate.com/golang/cook/connector"
	cook_sql "gitlab.niceprivate.com/golang/cook/sql/mysql"
	"testing"
)

type M_User struct {
	Id          int
	Name        string
	Password    string `orm-name:"passwd"`
	Avatar      string
	Description string
}

func F_User() interface{} {
	return &M_User{}
}

func init() {
	configs := map[string]cook_conn.MysqlConf{
		"default": cook_conn.MysqlConf{
			Addr:     "bdb.slave.kkgoo.niceprivate.com:3306",
			Username: "nice",
			Password: "Cb84eZaa229ddnm",
			Database: "kkgoo",

			MaxIdle: 4,
			MaxOpen: 4,
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
		cook_sql.E_eq("id", 5012470),
	)); err != nil {
		t.Logf("query faield: %s", err)
		t.Fail()
	}
	t.Logf("%#v", user)
	t.Fail()
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
