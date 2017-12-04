package orm

import (
	"database/sql"
	"errors"
	cook_conn "gitlab.niceprivate.com/golang/cook/connector"
	cook_opt "gitlab.niceprivate.com/golang/cook/option"
	cook_util "gitlab.niceprivate.com/golang/cook/util"
	"reflect"
	"strings"
)

var (
	ErrInvalidModelFactory        = errors.New("model's factory must return a pointer to struct")
	Heaven                 []*God = make([]*God, 0)
)

type FieldInfo struct {
	Model_field_name string
	Model_field      reflect.StructField
	DB_field_name    string
	Tags             map[string][]string
}
type ModelInfo struct {
	Pkg      string
	Name     string
	FullName string

	Model_field_mapping map[string]*FieldInfo
	DB_field_mapping    map[string]*FieldInfo
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
	db      *sql.DB            // handler
	Model   *ModelInfo
	Opts    *GodOptions
}

func (g *God) DB() *sql.DB {
	if g.db == nil {
		if db, err := cook_conn.GetMysql(g.Node); err != nil {
			return nil
		} else {
			g.db = db
		}
	}
	return g.db
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
		m_info *ModelInfo
		opts   *GodOptions
		err    error
	)

	if m_info, err = buildModelInfo(factory); err != nil {
		return nil, err
	}

	opts = &GodOptions{}
	cook_opt.Apply(opts, o...)

	god := &God{
		Factory: factory,
		Node:    node,
		Opts:    opts,
		Model:   m_info,
	}

	Heaven = append(Heaven, god)

	return god, nil
}

func buildModelInfo(factory func() interface{}) (*ModelInfo, error) {
	var (
		m         interface{}
		r_m       reflect.Type
		r_f       reflect.StructField
		m_info    *ModelInfo
		f_info    *FieldInfo
		db_f_name string
		tags      map[string][]string
	)
	m = factory()
	if r_m = reflect.TypeOf(m); r_m.Kind() != reflect.Ptr {
		return nil, ErrInvalidModelFactory
	}
	if r_m = r_m.Elem(); r_m.Kind() != reflect.Struct {
		return nil, ErrInvalidModelFactory
	}

	m_info = &ModelInfo{
		Pkg:                 r_m.PkgPath(),
		Name:                r_m.Name(),
		FullName:            r_m.PkgPath() + "." + r_m.Name(),
		Model_field_mapping: make(map[string]*FieldInfo),
		DB_field_mapping:    make(map[string]*FieldInfo),
	}
	for i := 0; i < r_m.NumField(); i++ {
		r_f = r_m.Field(i)

		tags = parseOrmTag(r_f)
		if args, ok := tags["col"]; ok && len(args) > 0 && len(args[0]) > 0 {
			db_f_name = args[0]
		} else {
			db_f_name = cook_util.Hump_to_underline(r_f.Name)
		}

		f_info = &FieldInfo{
			Model_field:      r_f,
			Model_field_name: r_f.Name,
			DB_field_name:    db_f_name,
			Tags:             tags,
		}
		m_info.Model_field_mapping[r_f.Name] = f_info
		m_info.DB_field_mapping[db_f_name] = f_info
	}
	return m_info, nil
}

// F(xxx,xxx);F;F(xxx)
func parseOrmTag(r_field reflect.StructField) map[string][]string {
	var (
		l_bracket_idx, r_bracket_idx int
		tags                         []string
		parsed                       map[string][]string = make(map[string][]string)
		tag_str                      string
	)
	tag_str = strings.TrimSpace(r_field.Tag.Get("orm"))
	tags = strings.Split(tag_str, ";")
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		l_bracket_idx = strings.Index(tag, "(")
		switch l_bracket_idx {
		case -1:
			parsed[tag] = []string{} // have no args
		case 0:
			continue // wrong format
		default:
			r_bracket_idx = strings.LastIndex(tag, ")")
			if r_bracket_idx <= l_bracket_idx {
				continue // wrong format
			}
			parsed[tag[:l_bracket_idx]] = strings.Split(tag[l_bracket_idx+1:r_bracket_idx], ",")
		}
	}
	return parsed
}
