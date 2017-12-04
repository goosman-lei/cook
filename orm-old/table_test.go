package orm

import (
	cook_conn "gitlab.niceprivate.com/golang/cook/connector"
	cook_sql "gitlab.niceprivate.com/golang/cook/sql/mysql"
	"testing"
	"time"
)

type M_Obj struct {
	Id     int
	Name   string
	Passwd string
}

func F_Obj() interface{} {
	return &M_Obj{}
}

func init() {
	configs := map[string]cook_conn.MysqlConf{
		"default": cook_conn.MysqlConf{
			Addr:     "bdb.slave.kkgoo.niceprivate.com:3306",
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

func Test_Table(t *testing.T) {
	var (
		god    *God
		eTable *cook_sql.Expr
		want   *cook_sql.Expr
	)

	god, _ = NewGod_shard_none(F_Obj, "default", "kk_user")
	eTable = god.Table()
	want = cook_sql.E_table("kk_user")
	if *eTable != *want {
		t.Logf("shard return: %#v, want: %#v", eTable, want)
		t.Fail()
	}

	god, _ = NewGod_shard_mod_int(F_Obj, "default", "kk_user_%d", 16)
	eTable = god.Table(38)
	want = cook_sql.E_table("kk_user_6")
	if *eTable != *want {
		t.Logf("shard return: %#v, want: %#v", eTable, want)
		t.Fail()
	}

	god, _ = NewGod_shard_mod_string_hex_8_to_int(F_Obj, "default", "kk_user_%d", 16)
	eTable = god.Table("Jack")
	want = cook_sql.E_table("kk_user_2")
	if *eTable != *want {
		t.Logf("shard return: %#v, want: %#v", eTable, want)
		t.Fail()
	}

}

func Test_Tables(t *testing.T) {
	var (
		god        *God
		eTables    []*cook_sql.Expr
		eArgs      [][][]interface{}
		wantTables []*cook_sql.Expr
		wantArgs   [][][]interface{}
	)

	god, _ = NewGod_shard_mod_int(F_Obj, "default", "kk_user_%d", 4)
	eTables, eArgs, _ = god.Tables([]int{3, 8, 7, 6, 2, 9, 1, 5, 13}, []string{"Jack", "Tom", "Ben", "Green", "Lucas", "John", "Bob", "Larunce", "Mark"})
	wantTables = []*cook_sql.Expr{
		cook_sql.E_table("kk_user_3"),
		cook_sql.E_table("kk_user_0"),
		cook_sql.E_table("kk_user_2"),
		cook_sql.E_table("kk_user_1"),
	}
	wantArgs = [][][]interface{}{
		[][]interface{}{
			[]interface{}{3, 7},
			[]interface{}{"Jack", "Ben"},
		},
		[][]interface{}{
			[]interface{}{8},
			[]interface{}{"Tom"},
		},
		[][]interface{}{
			[]interface{}{6, 2},
			[]interface{}{"Green", "Lucas"},
		},
		[][]interface{}{
			[]interface{}{9, 1, 5, 13},
			[]interface{}{"John", "Bob", "Larunce", "Mark"},
		},
	}

	for i, eTable := range eTables {
		if *eTable == *wantTables[0] {
			for j, wantCols := range wantArgs[0] {
				for k, col := range wantCols {
					if col != eArgs[i][j][k] {
						t.Logf("args[%d][%d] return: %#v, want: %#v", i, j, eArgs[i][j], wantCols)
						t.Fail()
					}
				}
			}
		} else if *eTable == *wantTables[1] {
			for j, wantCols := range wantArgs[1] {
				for k, col := range wantCols {
					if col != eArgs[i][j][k] {
						t.Logf("args[%d][%d] return: %#v, want: %#v", i, j, eArgs[i][j], wantCols)
						t.Fail()
					}
				}
			}
		} else if *eTable == *wantTables[2] {
			for j, wantCols := range wantArgs[2] {
				for k, col := range wantCols {
					if col != eArgs[i][j][k] {
						t.Logf("args[%d][%d] return: %#v, want: %#v", i, j, eArgs[i][j], wantCols)
						t.Fail()
					}
				}
			}
		} else if *eTable == *wantTables[3] {
			for j, wantCols := range wantArgs[3] {
				for k, col := range wantCols {
					if col != eArgs[i][j][k] {
						t.Logf("args[%d][%d] return: %#v, want: %#v", i, j, eArgs[i][j], wantCols)
						t.Fail()
					}
				}
			}
		} else {
			t.Logf("shard return unexpected table: %#v, args: %#v", eTable, eArgs[i])
			t.Fail()
		}
	}
}

func Benchmark_Tables(b *testing.B) {
	god, _ := NewGod_shard_mod_int(F_Obj, "default", "kk_user_%d", 4)

	for i := 0; i < b.N; i++ {
		god.Tables([]int{3, 8, 7, 6, 2, 9, 1, 5, 13}, []string{"Jack", "Tom", "Ben", "Green", "Lucas", "John", "Bob", "Larunce", "Mark"})
	}
}
