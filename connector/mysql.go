package connector

import (
	cook_log "cook/log"
	cook_util "cook/util"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type MysqlConnWrapper struct {
	*sql.DB
}

type MysqlConf struct {
	Addr     string
	Username string
	Password string
	Database string

	MaxIdle int
	MaxOpen int
}

var mysqlConnMapping *cook_util.CMap

func SetupMysql(configs map[string]MysqlConf) error {
	mysqlConnMapping = cook_util.NewCMap()
	for sn, config := range configs {
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
			config.Username, config.Password, config.Addr, config.Database))
		if err != nil {
			cook_log.Fatalf("mysql cluster [%s: %s@%s/%s] setup failed: %s",
				sn, config.Username, config.Addr, config.Database)
			return fmt.Errorf("open failed[%s]: %s@%s/%s", sn, config.Username, config.Addr, config.Database)
		}
		if err = db.Ping(); err != nil {
			cook_log.Fatalf("mysql cluster [%s: %s@%s/%s] ping failed when setup: %s",
				sn, config.Username, config.Addr, config.Database)
			return fmt.Errorf("ping failed[%s]: %s@%s/%s", sn, config.Username, config.Addr, config.Database)
		}
		db.SetMaxIdleConns(config.MaxIdle)
		db.SetMaxOpenConns(config.MaxOpen)
		mysqlConnMapping.Set(sn, db)
	}
	return nil
}

func GetMysql(sn string) (*MysqlConnWrapper, error) {
	if conn, exists := mysqlConnMapping.Get(sn); exists {
		return &MysqlConnWrapper{conn.(*sql.DB)}, nil
	}
	cook_log.Warnf("get mysql cluster[%s], but not ready", sn)
	return nil, fmt.Errorf("have no mysql cluster: %s", sn)
}

func MustGetMysql(sn string) *MysqlConnWrapper {
	conn, err := GetMysql(sn)
	if err != nil {
		panic(err)
	}
	return conn
}

func (m *MysqlConnWrapper) GetRow_Simple(table string, fields []string, where []MExpr) (*sql.Row, error) {
	return m.GetRow_Simple_Limit(table, fields, where, -1, -1)
}
func (m *MysqlConnWrapper) GetRow_Simple_Limit(table string, fields []string, where []MExpr, limit, offset int) (*sql.Row, error) {
	return m.GetRow_MTable_Limit(MT_Slice_Name(table), fields, where, limit, offset)
}
func (m *MysqlConnWrapper) GetRow_MTable(tables []MTable, fields []string, where []MExpr) (*sql.Row, error) {
	return m.GetRow_MTable_Limit(tables, fields, where, -1, -1)
}
func (m *MysqlConnWrapper) GetRow_MTable_Limit(tables []MTable, fields []string, where []MExpr, limit, offset int) (*sql.Row, error) {
	return m.GetRow(tables, ME_Slice_Fields(fields), where, limit, offset, EmptyMExprSlice, EmptyMExprSlice, EmptyMExprSlice)
}
func (m *MysqlConnWrapper) GetRow(tables []MTable, fields []MExpr, where []MExpr, limit, offset int, orderBy, groupBy, having []MExpr) (*sql.Row, error) {
	query, args, err := parseQuery(tables, fields, where, limit, offset, orderBy, groupBy, having)
	if err != nil {
		return nil, err
	}

	return m.QueryRow(query, args...), nil
}

func (m *MysqlConnWrapper) GetRows_Simple(table string, fields []string, where []MExpr) (*sql.Rows, error) {
	return m.GetRows_Simple_Limit(table, fields, where, -1, -1)
}
func (m *MysqlConnWrapper) GetRows_Simple_Limit(table string, fields []string, where []MExpr, limit, offset int) (*sql.Rows, error) {
	return m.GetRows_MTable_Limit(MT_Slice_Name(table), fields, where, limit, offset)
}
func (m *MysqlConnWrapper) GetRows_MTable(tables []MTable, fields []string, where []MExpr) (*sql.Rows, error) {
	return m.GetRows_MTable_Limit(tables, fields, where, -1, -1)
}
func (m *MysqlConnWrapper) GetRows_MTable_Limit(tables []MTable, fields []string, where []MExpr, limit, offset int) (*sql.Rows, error) {
	return m.GetRows(tables, ME_Slice_Fields(fields), where, limit, offset, EmptyMExprSlice, EmptyMExprSlice, EmptyMExprSlice)
}
func (m *MysqlConnWrapper) GetRows(tables []MTable, fields []MExpr, where []MExpr, limit, offset int, orderBy, groupBy, having []MExpr) (*sql.Rows, error) {
	query, args, err := parseQuery(tables, fields, where, limit, offset, orderBy, groupBy, having)
	if err != nil {
		return nil, err
	}

	return m.Query(query, args...)
}

func (m *MysqlConnWrapper) Count(tables []MTable, where []MExpr, countField string) (count int, err error) {
	row, err := m.GetRow(tables, []MExpr{
		ME_Literal("COUNT(" + countField + ") AS count"),
	}, where, -1, -1, EmptyMExprSlice, EmptyMExprSlice, EmptyMExprSlice)
	if err == nil {
		err = row.Scan(&count)
	}
	return
}

func (m *MysqlConnWrapper) Update(tables []MTable, sets []MExpr, where []MExpr, orderBy []MExpr, limit int) (nAffected int64, err error) {
	query, args, err := parseUpdate(tables, sets, where, orderBy, limit)
	if err != nil {
		return
	}

	result, err := m.Exec(query, args)
	if err == nil {
		nAffected, err = result.RowsAffected()
	}
	return
}

func (m *MysqlConnWrapper) Insert(table string, sets []MExpr, onDup []MExpr) (lastId int64, err error) {
	query, args, err := parseInsert(table, sets, onDup)
	if err != nil {
		return
	}

	result, err := m.Exec(query, args)
	if err == nil {
		lastId, err = result.LastInsertId()
	}
	return
}

func (m *MysqlConnWrapper) Delete(table string, where []MExpr, orderBy []MExpr, limit int) (nAffected int64, err error) {
	query, args, err := parseDelete(table, where, orderBy, limit)
	if err != nil {
		return
	}

	result, err := m.Exec(query, args)
	if err == nil {
		nAffected, err = result.RowsAffected()
	}
	return

}
