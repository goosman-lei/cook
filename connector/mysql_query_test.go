package connector

import (
	"testing"
)

var (
	bTables  []MTable
	bFields  []MExpr
	bWhere   []MExpr
	bLimit   int
	bOffset  int
	bGroupBy []MExpr
	bOrderBy []MExpr
	bHaving  []MExpr
)

func init() {
	bTables = []MTable{
		MTable{name: "kk_user", alias: "u", idx_hint_list: []string{"force idx_uname"}},
		MTable{join: MYSQL_JOIN_LEFT, name: "kk_user_bind_account", alias: "uba", on: ME_Literal("u.id = uba.uid")},
	}
	bFields = []MExpr{
		ME_Field("id"),
		ME_Field("name"),
		ME_Field("age"),
	}
	bWhere = []MExpr{
		ME_Ne("score", 80),
		ME_Gt("score", 80),
		ME_Ge("score", 80),
		ME_Eq("score", 80),
		ME_Le("score", 80),
		ME_Lt("score", 80),
		ME_Between("score", 80, 100),
		ME_Lt("score", 8.3),
		ME_Eq("name", "Jack"),
		ME_Like("name", "Jack%"),
		ME_Eq("pass", true),
		ME_Or([]MExpr{
			ME_NotLike("name", "Jack_"),
			ME_In("sex", []interface{}{"man", "woman"}),
			ME_NotIn("color", []interface{}{"blue", "black"}),
		}),
		ME_Not(ME_NotNull("money")),
		ME_And([]MExpr{
			ME_IsNull("color"),
		}),
		ME_Literal("uid = 100"),
	}
	bLimit = 10
	bOffset = 100
	bGroupBy = []MExpr{
		ME_Field("name").Asc(),
		ME_Alias("pass", "pass_alias").Desc(),
	}
	bOrderBy = []MExpr{
		ME_Field("score").Desc(),
		ME_Literal("age * age").Desc(),
	}
	bHaving = []MExpr{
		ME_Gt("score", 30),
		ME_Between("score", 80, 100),
		ME_Lt("score", 8.3),
		ME_Eq("name", "Tom"),
	}

}

func TestParseComplexQuery(t *testing.T) {
	query, args, err := parseQuery([]MTable{
		MTable{name: "kk_user", alias: "u", idx_hint_list: []string{"force idx_uname"}},
		MTable{join: MYSQL_JOIN_LEFT, name: "kk_user_bind_account", alias: "uba", on: ME_Literal("u.id = uba.uid")},
	}, []MExpr{
		ME_Field("id"),
		ME_Field("name"),
		ME_Field("age"),
	}, []MExpr{
		ME_Ne("score", 80),
		ME_Gt("score", 80),
		ME_Ge("score", 80),
		ME_Eq("score", 80),
		ME_Le("score", 80),
		ME_Lt("score", 80),
		ME_Between("score", 80, 100),
		ME_Lt("score", 8.3),
		ME_Eq("name", "Jack"),
		ME_Like("name", "Jack%"),
		ME_Eq("pass", true),
		ME_Or([]MExpr{
			ME_NotLike("name", "Jack_"),
			ME_In("sex", []interface{}{"man", "woman"}),
			ME_NotIn("color", []interface{}{"blue", "black"}),
		}),
		ME_Not(ME_NotNull("money")),
		ME_And([]MExpr{
			ME_IsNull("color"),
		}),
		ME_Literal("uid = 100"),
	}, 10, 100, []MExpr{
		ME_Field("name").Asc(),
		ME_Alias("pass", "pass_alias").Desc(),
	}, []MExpr{
		ME_Field("score").Desc(),
		ME_Literal("age * age").Desc(),
	}, []MExpr{
		ME_Gt("score", 30),
		ME_Between("score", 80, 100),
		ME_Lt("score", 8.3),
		ME_Eq("name", "Tom"),
	})

	if err != nil {
		t.Logf("parsed SQL unexpected error: %s", err)
		t.Fail()
		return
	}

	if query != "SELECT id, name, age FROM kk_user AS u force idx_uname LEFT JOIN kk_user_bind_account AS uba ON u.id = uba.uid WHERE score != ? AND score > ? AND score >= ? AND score = ? AND score <= ? AND score < ? AND score BETWEEN ? AND ? AND score < ? AND name = ? AND name LIKE ? AND pass = ? AND (name NOT LIKE ? OR sex IN(?,?) OR color NOT IN(?,?)) AND NOT (money IS NOT NULL) AND (color IS NULL) AND uid = 100 GROUP BY score DESC, age * age DESC HAVING score > ? AND score BETWEEN ? AND ? AND score < ? AND name = ? ORDER BY name ASC, pass AS pass_alias DESC LIMIT 10 OFFSET 100" {
		t.Logf("parsed SQL unexpected: %s", query)
		t.Fail()
	}

	for i := 0; i < 7; i++ {
		if args[i].(int) != 80 {
			t.Logf("args[%d] != 80", i)
			t.Fail()
		}
	}

	if args[7].(int) != 100 {
		t.Logf("args[7] != 100")
		t.Fail()
	}

	if args[8].(float64) != 8.3 {
		t.Logf("args[8] != 8.3")
		t.Fail()
	}

	if args[9].(string) != "Jack" {
		t.Logf("args[9] != Jack")
		t.Fail()
	}

	if args[11].(bool) != true {
		t.Logf("args[11] != true")
		t.Fail()
	}

	if args[15].(string) != "blue" {
		t.Logf("args[15] != blue")
		t.Fail()
	}

	if len(args) != 22 {
		t.Logf("len(args) != 22")
		t.Fail()
	}
}

func TestParseUpdate(t *testing.T) {
	query, args, err := parseUpdate([]MTable{
		MTable{name: "kk_user", alias: "u"},
		MTable{name: "kk_user_bind_account", alias: "uba", join: MYSQL_JOIN_LEFT, on: ME_Literal("uba.uid = u.id")},
	}, []MExpr{
		ME_Assign("u.age", ME_Literal("u.age + 1")),
		ME_Assign("uba.mobile", ME_Literal("u.wid")),
	}, []MExpr{
		ME_Gt("uba.age", 10),
		ME_Lt("u.age", 100),
	}, []MExpr{
		ME_Field("uba.age").Desc(),
		ME_Field("u.id").Asc(),
	}, 10)

	if err != nil {
		t.Logf("parsed SQL unexpected error: %s", err)
		t.Fail()
		return
	}

	if query != "UPDATE kk_user AS u LEFT JOIN kk_user_bind_account AS uba ON uba.uid = u.id SET u.age = (u.age + 1), uba.mobile = (u.wid) WHERE uba.age > ? AND u.age < ? ORDER BY uba.age DESC, u.id ASC LIMIT 10" {
		t.Logf("parsed SQL unexpected: %s", query)
		t.Fail()
	}

	if args[0].(int) != 10 {
		t.Logf("args[0] != 10")
		t.Fail()
	}

	if len(args) != 2 {
		t.Logf("len(args) != 2")
		t.Fail()
	}
}

func TestParseInsert(t *testing.T) {
	query, args, err := parseInsert("kk_user", []MExpr{
		ME_Assign_Value("name", "goosman"),
		ME_Assign_Value("mobile", 15210643724),
		ME_Assign_Value("country", 1),
		ME_Assign_Value("height", 1.65),
	}, []MExpr{
		ME_Literal("height = height + 0.5"),
		ME_Literal("age = age - 10"),
	})

	if err != nil {
		t.Logf("parsed SQL unexpected error: %s", err)
		t.Fail()
		return
	}

	if query != "INSERT INTO kk_user SET name = ?, mobile = ?, country = ?, height = ? ON DUPLICATE KEY UPDATE height = height + 0.5, age = age - 10" {
		t.Logf("parsed SQL unexpected: %s", query)
		t.Fail()
	}

	if args[0].(string) != "goosman" {
		t.Logf("args[0] != goosman")
		t.Fail()
	}

	if len(args) != 4 {
		t.Logf("len(args) != 4")
		t.Fail()
	}
}

func TestParseDelete(t *testing.T) {
	query, args, err := parseDelete("kk_user", []MExpr{
		ME_Ge("height", 1.5),
		ME_Between("age", 25, 38),
	}, []MExpr{
		ME_Field("age").Desc(),
		ME_Expr_Alias(ME_Literal("age + height * 10"), "age_height").Asc(),
	}, 10)

	if err != nil {
		t.Logf("parsed SQL unexpected error: %s", err)
		t.Fail()
		return
	}

	if query != "DELETE FROM kk_user WHERE height >= ? AND age BETWEEN ? AND ? ORDER BY age DESC, (age + height * 10) AS age_height ASC LIMIT 10" {
		t.Logf("parsed SQL unexpected: %s", query)
		t.Fail()
	}

	if args[1].(int) != 25 {
		t.Logf("args[1] != 25")
		t.Fail()
	}

	if len(args) != 3 {
		t.Logf("len(args) != 3")
		t.Fail()
	}
}

func BenchmarkParseComplex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseQuery(bTables, bFields, bWhere, bLimit, bOffset, bGroupBy, bOrderBy, bHaving)
	}
}

func BenchmarkParseSimple_1(b *testing.B) {
	sTables := []MTable{
		MTable{name: "kk_user"},
	}
	sFields := []MExpr{
		ME_Field("id"),
		ME_Field("name"),
	}
	sWhere := []MExpr{
		ME_In("name", []interface{}{"ALEX", "大鹏", "bluef"}),
		ME_Gt("age", 20),
	}
	sLimit := 10
	sOffset := 0
	sGroupBy := []MExpr{}
	sOrderBy := []MExpr{
		ME_Field("id").Desc(),
	}
	sHaving := []MExpr{}

	for i := 0; i < b.N; i++ {
		parseQuery(sTables, sFields, sWhere, sLimit, sOffset, sGroupBy, sOrderBy, sHaving)
	}
}

func BenchmarkParseSimple_2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sTables := []MTable{
			MTable{name: "kk_user"},
		}
		sFields := []MExpr{
			ME_Field("id"),
			ME_Field("name"),
		}
		sWhere := []MExpr{
			ME_In("name", []interface{}{"ALEX", "大鹏", "bluef"}),
			ME_Gt("age", 20),
		}
		sLimit := 10
		sOffset := 0
		sGroupBy := []MExpr{}
		sOrderBy := []MExpr{
			ME_Field("id").Desc(),
		}
		sHaving := []MExpr{}

		parseQuery(sTables, sFields, sWhere, sLimit, sOffset, sGroupBy, sOrderBy, sHaving)
	}
}
