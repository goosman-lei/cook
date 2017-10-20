package model

import (
	"database/sql"
	"errors"
	cook_conn "gitlab.niceprivate.com/golang/cook/connector"
	cook_opt "gitlab.niceprivate.com/golang/cook/option"
	cook_util "gitlab.niceprivate.com/golang/cook/util"
	"reflect"
)

var (
	ErrInvalidModelFactory = errors.New("model's factory must return a pointer to struct")
)

type FieldInfo struct {
	Name    string
	Field   reflect.StructField
	DB_name string
}

type ModelInfo struct {
	FieldMapping   map[string]*FieldInfo
	DBFieldMapping map[string]*FieldInfo
}

type GodOptions struct {
	TableFmt  string
	ShardType uint8

	ShardCnt int
	HashAlgo string
}

type God struct {
	Factory func() interface{} // factory of model
	Node    string             // used node
	DB      *sql.DB            // handler
	Model   ModelInfo
	Opts    *GodOptions
}

func NewGod_shard_none(factory func() interface{}, node string, name string) (*God, error) {
	return newGod(
		factory,
		node,
		cook_opt.Opt("TableFmt", name),
		cook_opt.Opt("ShardType", SHARD_TYPE_NONE),
	)
}

func NewGod_shard_mod_int(factory func() interface{}, node string, format string, shard_count int) (*God, error) {
	return newGod(
		factory,
		node,
		cook_opt.Opt("TableFmt", format),
		cook_opt.Opt("ShardType", SHARD_TYPE_MOD_INT),
		cook_opt.Opt("ShardCnt", shard_count),
	)
}

func NewGod_shard_mod_string_hex_8_to_int(factory func() interface{}, node string, format string, shard_count int) (*God, error) {
	return newGod(
		factory,
		node,
		cook_opt.Opt("TableFmt", format),
		cook_opt.Opt("ShardType", SHARD_TYPE_MOD_STRING),
		cook_opt.Opt("HashAlgo", "hex_0_8_to_int"),
		cook_opt.Opt("ShardCnt", shard_count),
	)
}

func newGod(factory func() interface{}, node string, o ...cook_opt.Option) (*God, error) {
	var (
		m         interface{}  // model object
		r_m       reflect.Type // reflect of model object's type
		m_info    ModelInfo
		r_f       reflect.StructField // reflect of model objcet's field(type)
		f_i       *FieldInfo
		db_f_name string // field name in db
		ok        bool   // for check model's tag
	)
	// build model info
	m = factory()
	if r_m = reflect.TypeOf(m); r_m.Kind() != reflect.Ptr {
		return nil, ErrInvalidModelFactory
	}
	if r_m = r_m.Elem(); r_m.Kind() != reflect.Struct {
		return nil, ErrInvalidModelFactory
	}
	m_info = ModelInfo{
		FieldMapping:   make(map[string]*FieldInfo),
		DBFieldMapping: make(map[string]*FieldInfo),
	}
	for i := 0; i < r_m.NumField(); i++ {
		r_f = r_m.Field(i)
		if db_f_name, ok = r_f.Tag.Lookup("orm-name"); !ok {
			db_f_name = cook_util.Hump_to_underline(r_f.Name)
		}
		f_i = &FieldInfo{
			Name:    r_f.Name,
			Field:   r_f,
			DB_name: db_f_name,
		}
		m_info.FieldMapping[r_f.Name] = f_i
		m_info.DBFieldMapping[db_f_name] = f_i
	}

	opts := &GodOptions{}
	cook_opt.Apply(opts, o...)
	if db, err := cook_conn.GetMysql(node); err != nil {
		return nil, err
	} else {
		return &God{
			Factory: factory,
			Node:    node,
			Opts:    opts,
			DB:      db,
			Model:   m_info,
		}, nil
	}
}
