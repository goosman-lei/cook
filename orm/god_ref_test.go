package orm

import (
	"gitlab.niceprivate.com/golang/cook/orm/table"
	"testing"
)

type M_User_GodRefTest struct {
	*M       `orm:"nomapping"`
	Id       int `orm:"pk"`
	Name     string
	Password string `orm:"col(passwd)"`
}

func Test_God_Ref(t *testing.T) {
	GodOf_User := NewGod((*M_User_GodRefTest)(nil), "default", table.Table_normal("kk_user"))

	if GodOf_User.Model == nil {
		t.Logf("model reflect failed")
		t.Fail()
	}

	if len(GodOf_User.Model.PK) != 1 || GodOf_User.Model.PK[0].R_StructField.Name != "Id" {
		t.Logf("PK reflect failed")
		t.Fail()
	}

	if len(GodOf_User.Model.Mapping_with_field) != 3 {
		t.Logf("Mapping_with_field reflect failed")
		t.Fail()
	}

	mapping_field := GodOf_User.Model.Mapping_with_field
	mapping_column := GodOf_User.Model.Mapping_with_column
	if mapping_field["Id"].Column != "id" || mapping_column["id"].R_StructField.Name != "Id" {
		t.Logf("mapping reflect failed")
		t.Fail()
	}
	if mapping_field["Name"].Column != "name" || mapping_column["name"].R_StructField.Name != "Name" {
		t.Logf("mapping reflect failed")
		t.Fail()
	}
	if mapping_field["Password"].Column != "passwd" || mapping_column["passwd"].R_StructField.Name != "Password" {
		t.Logf("mapping reflect failed")
		t.Fail()
	}
}
