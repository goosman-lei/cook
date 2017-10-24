package orm

import (
	"database/sql"
	"net"
	"testing"
	"time"
)

var skipMysql bool = false

type FieldDesc struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default string
	Extra   string
}

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

			MaxIdle: 4,
			MaxOpen: 4,
		},
		"cluster2": MysqlConf{
			Addr:     "bdb.slave.kkgoo.niceprivate.com:3306",
			Username: "nice",
			Password: "Cb84eZaa229ddnm",
			Database: "kkgoo",

			MaxIdle: 4,
			MaxOpen: 4,
		},
	}

	SetupMysql(configs)
}

func Test_pool(t *testing.T) {
	if skipMysql {
		return
	}
	var (
		db     *sql.DB
		rows   *sql.Rows
		fields []*FieldDesc
	)
	db, _ = GetMysql("cluster1")
	rows, _ = db.Query("DESC kk_user")

	for rows.Next() {
		field_desc := new(FieldDesc)
		rows.Scan(&field_desc.Field, &field_desc.Type, &field_desc.Null, &field_desc.Key, &field_desc.Default, &field_desc.Extra)
		fields = append(fields, field_desc)
	}

	t.Logf("+------------------+--------------------------------+------+-----+---------+----------------+")
	t.Logf("| Field            | Type                           | Null | Key | Default | Extra          |")
	t.Logf("+------------------+--------------------------------+------+-----+---------+----------------+")
	for _, field := range fields {
		t.Logf("| %-16s | %-30s | %-4s | %-3s | %-7s | %-14s |\n", field.Field, field.Type, field.Null, field.Key, field.Default, field.Extra)
	}
	t.Logf("+------------------+--------------------------------+------+-----+---------+----------------+")

	//t.Fail()
}
