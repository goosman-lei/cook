package orm

import (
	"gitlab.niceprivate.com/golang/cook/orm-ng/table"
	"testing"
)

type M_User struct {
	Id       int `orm:"pk"`
	Name     string
	Password string `orm:"col(passwd)"`
}

func Test_Load(t *testing.T) {
	GodOf_User := NewGod(&M_User{}, "default", table.Table_normal("kk_user"))

	if GodOf_User.Model == nil {
		t.Logf("model reflect failed")
		t.Fail()
	}

	if len(GodOf_User.Model.PK) != 1 || GodOf_User.Model.PK[0].R_StructField.Name != "Id" {
		t.Logf("PK reflect failed")
		t.Fail()
	}

	if len(GodOf_User.Model.Fields) != 3 {
		t.Logf("Fields reflect failed")
		t.Fail()
	}

	fields := GodOf_User.Model.Fields
	if fields[0].R_StructField.Name != "Id" || fields[0].Column != "id" {
		t.Logf("Fields reflect failed")
		t.Fail()
	}
	if fields[1].R_StructField.Name != "Name" || fields[1].Column != "name" {
		t.Logf("Fields reflect failed")
		t.Fail()
	}
	if fields[2].R_StructField.Name != "Password" || fields[2].Column != "passwd" {
		t.Logf("Fields reflect failed")
		t.Fail()
	}
}
