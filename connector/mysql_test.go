package connector

import (
	"database/sql"
	"net"
	"testing"
	"time"
)

var skipMysql bool = false

func init() {
	testAddr := "bdb.slave.kkgoo.niceprivate.com:3306"
	if _, err := net.DialTimeout("tcp", testAddr, 100*time.Millisecond); err != nil {
		skipMysql = true
		return
	}

	configs := map[string]MysqlConf{
		"cluster1": MysqlConf{
			Addr:     "bdb.slave.kkgoo.niceprivate.com:3306",
			Username: "nice",
			Password: "Cb84eZaa229ddnm",
			Database: "kkgoo",

			MaxIdle: 16,
			MaxOpen: 16,
		},
		"cluster2": MysqlConf{
			Addr:     "bdb.slave.kkgoo.niceprivate.com:3306",
			Username: "nice",
			Password: "Cb84eZaa229ddnm",
			Database: "kkgoo",

			MaxIdle: 16,
			MaxOpen: 16,
		},
	}

	SetupMysql(configs)

}

func TestQuery(t *testing.T) {
	if skipMysql {
		t.Skipf("test mysql server is not reachable")
	}
	var (
		id          int
		name        string
		conn1       = MustGetMysql("cluster1")
		queryMUinfo = "SELECT id, name FROM kk_user WHERE name IN(?, ?, ?, ?)"
	)

	rows, err := conn1.Query(queryMUinfo, "bluef", "狮你妹", "大鹏", "ALEX")
	if err != nil {
		t.Logf("Query[%s] failed: %s", queryMUinfo, err)
		t.Fail()
	}
	for rows.Next() {
		rows.Scan(&id, &name)
		switch id {
		case 116074:
			if name != "bluef" {
				t.Logf("unexptected data: uid = %d uname = %s", id, name)
				t.Fail()
			}
		case 93:
			if name != "大鹏" {
				t.Logf("unexptected data: uid = %d uname = %s", id, name)
				t.Fail()
			}
		case 94:
			if name != "ALEX" {
				t.Logf("unexptected data: uid = %d uname = %s", id, name)
				t.Fail()
			}
		case 100:
			if name != "狮你妹" {
				t.Logf("unexptected data: uid = %d uname = %s", id, name)
				t.Fail()
			}
		default:
			t.Logf("unexptected data: uid = %d uname = %s", id, name)
			t.Fail()
		}
	}
	rows.Close()
}

func TestGetRow(t *testing.T) {
	if skipMysql {
		t.Skipf("test mysql server is not reachable")
	}

	row, err := MustGetMysql("cluster1").GetRow_Simple("kk_user", []string{"id", "name"}, []MExpr{
		ME_Eq("id", 0xFFFFFFFF),
	})
	if err != nil {
		t.Logf("GetRow failed: %s", err)
		t.Fail()
	}
	var (
		id   int
		name string
	)
	if err = row.Scan(&id, &name); err != sql.ErrNoRows {
		t.Logf("want error: %s, actual: %s", sql.ErrNoRows, err)
		t.Fail()
	}
}

func TestGetRows(t *testing.T) {
	if skipMysql {
		t.Skipf("test mysql server is not reachable")
	}
	var (
		id    int
		name  string
		conn1 = MustGetMysql("cluster1")
		rows  *sql.Rows
		err   error
	)
	rows, err = conn1.GetRows_Simple("kk_user", []string{"id", "name"}, []MExpr{
		ME_Or([]MExpr{
			ME_In("name", []interface{}{"大鹏", "狮你妹", "ALEX", "bluef"}),
			ME_Eq("id", 5012470),
		}),
	})
	if err != nil {
		t.Logf("Query failed: %s", err)
		t.Fail()
	}
	for rows.Next() {
		rows.Scan(&id, &name)
		switch id {
		case 116074:
			if name != "bluef" {
				t.Logf("unexptected data: uid = %d uname = %s", id, name)
				t.Fail()
			}
		case 93:
			if name != "大鹏" {
				t.Logf("unexptected data: uid = %d uname = %s", id, name)
				t.Fail()
			}
		case 94:
			if name != "ALEX" {
				t.Logf("unexptected data: uid = %d uname = %s", id, name)
				t.Fail()
			}
		case 100:
			if name != "狮你妹" {
				t.Logf("unexptected data: uid = %d uname = %s", id, name)
				t.Fail()
			}
		case 5012470:
			if name != "今晚打虫子" {
				t.Logf("unexptected data: uid = %d uname = %s", id, name)
				t.Fail()
			}
		default:
			t.Logf("unexptected data: uid = %d uname = %s", id, name)
			t.Fail()
		}
	}

}

func TestCount(t *testing.T) {
	if skipMysql {
		t.Skipf("test mysql server is not reachable")
	}
	var (
		conn1 = MustGetMysql("cluster1")
		count int
		err   error
	)
	count, err = conn1.Count(MT_Slice_Name("kk_user"), []MExpr{
		ME_Or([]MExpr{
			ME_In("id", []interface{}{93, 94, 100}),
			ME_Eq("id", 5012470),
		}),
	}, "id")
	if err != nil {
		t.Logf("Query failed: %s", err)
		t.Fail()
	}
	if count != 4 {
		t.Logf("Count failed: wanted = %d, actually = %d", 4, count)
		t.Fail()
	}
}
