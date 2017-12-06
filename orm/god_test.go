package orm

import (
	"testing"
)

type M_User_GodTest struct {
	*M       `orm:"nomapping"`
	Id       int `orm:"pk"`
	Name     string
	Password string `orm:"col(passwd)"`
}

func F_User_GodTest() Model {
	return &M_User_GodTest{}
}

func Test_Tpl(t *testing.T) {
	GodOf_User := NewGod(F_User_GodTest, "default", Table_normal("kk_user"))
	GodOf_User.Tpl("simple", "id", "name")
	GodOf_User.Tpl("full", "id", "name", "passwd", E_field("name").Alias("nick_name"))

	if _, ok := GodOf_User.Tpls["simple"]; !ok {
		t.Fatalf("register tpl simple failed")
	}
	if len(GodOf_User.Tpls["simple"]) != 2 {
		t.Fatalf("register tpl simple failed")
	}
	if GodOf_User.Tpls["simple"][0].op != OP_FIELD || GodOf_User.Tpls["simple"][0].field != "id" {
		t.Fatalf("register tpl simple failed")
	}
	if GodOf_User.Tpls["simple"][1].op != OP_FIELD || GodOf_User.Tpls["simple"][1].field != "name" {
		t.Fatalf("register tpl simple failed")
	}

	if _, ok := GodOf_User.Tpls["full"]; !ok {
		t.Fatalf("register tpl full failed")
	}
	if len(GodOf_User.Tpls["full"]) != 4 {
		t.Fatalf("register tpl full failed")
	}
	if GodOf_User.Tpls["full"][0].op != OP_FIELD || GodOf_User.Tpls["full"][0].field != "id" {
		t.Fatalf("register tpl full failed")
	}
	if GodOf_User.Tpls["full"][1].op != OP_FIELD || GodOf_User.Tpls["full"][1].field != "name" {
		t.Fatalf("register tpl full failed")
	}
	if GodOf_User.Tpls["full"][2].op != OP_FIELD || GodOf_User.Tpls["full"][2].field != "passwd" {
		t.Fatalf("register tpl full failed")
	}
	if GodOf_User.Tpls["full"][3].op != OP_FIELD || GodOf_User.Tpls["full"][3].field != "name" || GodOf_User.Tpls["full"][3].alias != "nick_name" {
		t.Fatalf("register tpl full failed")
	}
}
