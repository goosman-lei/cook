package orm

import (
	"testing"
)

type M_User_GodApiBenchmark struct {
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

func F_User_GodApiBenchmark() Model {
	return &M_User_GodApiBenchmark{}
}

var (
	GodOf_User_ApiBenchmark *God = NewGod(F_User_GodApiBenchmark, "not-exists-node", Table_normal("kk_user"))
)

func init() {
}

/*
3us/op
*/
func Benchmark_God_Multi(b *testing.B) {
	var users []*M_User_GodApiBenchmark

	for i := 0; i < b.N; i++ {
		GodOf_User_ApiBenchmark.On(E_eq("sex", "male")).Groupby("age").Having(E_gt("add_time", 0)).Orderby(E_desc("id")).Limit(10, 20).Multi(&users, "id", "name", E_literal("unix_timestamp() - add_time AS reg_secs"))
	}
	//b.Logf("%v", GodOf_User_ApiBenchmark.LastStatement.SQL)
	//b.Fail()
}

/*
40us/op
*/
func Benchmark_God_Insert(b *testing.B) {
	users := []interface{}{
		&M_User_GodApiBenchmark{Id: 1, Name: "Goosman-lei", Password: "HAHAHAHA", AddTime: 191919191, Sex: "male"},
		&M_User_GodApiBenchmark{Id: 2, Name: "goosman", Password: "HEHEHEHE", AddTime: 29292929, Sex: "male"},
		&M_User_GodApiBenchmark{Id: 3, Name: "goosman", Password: "HEHEHEHE", AddTime: 29292929, Sex: "male"},
		&M_User_GodApiBenchmark{Id: 4, Name: "goosman", Password: "HEHEHEHE", AddTime: 29292929, Sex: "male"},
		&M_User_GodApiBenchmark{Id: 5, Name: "goosman", Password: "HEHEHEHE", AddTime: 29292929, Sex: "male"},
		&M_User_GodApiBenchmark{Id: 6, Name: "goosman", Password: "HEHEHEHE", AddTime: 29292929, Sex: "male"},
		&M_User_GodApiBenchmark{Id: 7, Name: "goosman", Password: "HEHEHEHE", AddTime: 29292929, Sex: "male"},
		&M_User_GodApiBenchmark{Id: 8, Name: "goosman", Password: "HEHEHEHE", AddTime: 29292929, Sex: "male"},
		&M_User_GodApiBenchmark{Id: 9, Name: "goosman", Password: "HEHEHEHE", AddTime: 29292929, Sex: "male"},
		&M_User_GodApiBenchmark{Id: 10, Name: "goosman", Password: "HEHEHEHE", AddTime: 29292929, Sex: "male"},
		[]string{"id", "name", "passwd", "add_time", "sex"},
	}

	for i := 0; i < b.N; i++ {
		GodOf_User_ApiBenchmark.Insert(users...)
	}
	//b.Logf("%v", GodOf_User_ApiBenchmark.LastStatement.Args)
	//b.Fail()
}

/*
1.5us/op
*/
func Benchmark_God_Delete(b *testing.B) {
	user := M_User_GodApiBenchmark{Id: 1}

	for i := 0; i < b.N; i++ {
		GodOf_User_ApiBenchmark.Delete(&user)
	}
	//b.Logf("%v", GodOf_User_ApiBenchmark.LastStatement.SQL)
	//b.Fail()
}

/*
6us/op
*/
func Benchmark_God_Update(b *testing.B) {
	user := M_User_GodApiBenchmark{Id: 2, Name: "goosman", Password: "HEHEHEHE", AddTime: 29292929, Sex: "male"}

	for i := 0; i < b.N; i++ {
		GodOf_User_ApiBenchmark.Update(&user, []string{"name", "passwd", "add_time", "sex"})
	}
	//b.Logf("%v", GodOf_User_ApiBenchmark.LastStatement.SQL)
	//b.Fail()
}
