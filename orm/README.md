# cook/orm

## 特性

* 主从分离

* 水平分表

* 对象映射

* 结构化查询

## Demo for common usage

### 按照下面规则定义你的Model

* 继承\*orm.M, 并标记为nomapping

* model中可用的tag

  * nomapping: 不和数据库字段映射绑定

  * pk: 主键. 用于Load/Loads方法

  * col(<name>): 定义对应的列名, 未指定则默认按照驼峰转下划线规则处理

* 定义Model的工厂方法, 返回空的Model指针

```
import (
    "gitlab.niceprivate.com/golang/cook/orm"
)

type M_User struct {
	*orm.M       `orm:"nomapping"`
	Id       int `orm:"pk"`
	Name     string
	Password string `orm:"col(passwd)"`
	Avatar   string
	Desc     string `orm:"col(description)"`
	AddTime  int
	Age      int
	Sex      string
}

func F_User() Model { return &M_User{} }
```

### 注册绑定God对象

给自己模块, 定义一个全局的God对象, 用来作为后续所有Model操作的入口

```
var (
// 普通God
    GodOf_User *God = NewGod(F_User, "数据库节点名", Table_normal("kk_user"))

// 主从分离God
    GodOf_User *God = NewGod_master_slave(F_User, "主库节点名", "从库节点名", Table_normal("kk_user"))

// 分表God
    GodOf_User *God = NewGod(F_User, "数据库节点名", Table_mod_int("kk_user_%d", 256))
)
```

### 让God带你起飞

```
// load with pk
user, err := GodOf_User.Load(1)
user.(*M_User).Name

// loads with pk
users, err := GodOf_User.Load(1, 2, 3, 4)
users[0].(*M_User).Name

// Count
cnt, err := GodOf_User.On(E_eq("age", 18)).Count()

// query one
user, err := GodOf_User.On(E_eq("id", 18)).One()

// query one. with optional columns
user, err := GodOf_User.On(E_eq("id", 18)).One("id", "name", E_literal("age + 1 AS plus_age"))
user.(*M_User).Id
user.(*M_User).Name
user.(*M_User).Int("plus_age")  // 查询数据库返回的列如果没有在model中设定映射关系, 则会存储到Model.M.E的map中, 可以使用Model的相关类型方法读取

// query multi. with optional columns
users, err := GodOf_User.On(E_in("id", []int{1, 2, 3, 4, 5})).One("id", "name", E_literal("age + 1 AS plus_age"))

// update
user = M_User{
    Id: 1,
    Name: "Goosman-lei",
}
affnum, err := GodOf_User.Update(&user)

// insert
affnum, err := GodOf_User.Insert(&user)

affnum, err := GodOf_User.Delete(&user)
affnum, err := GodOf_User.Delete(1)
affnum, err := GodOf_User.Delete(E_in("id", []{1, 2, 3, 4}))
```

## API

```
// E_xxx系列方法, 用来描述SQL语句中的表达式
func E_fields(fs ...string) []*Expr
type Expr struct{ ... }
    func E_and(d ...*Expr) *Expr
    func E_asc(f string) *Expr
    func E_assign(f string, d interface{}) *Expr
    func E_between(f string, d1, d2 interface{}) *Expr
    func E_comma(d ...*Expr) *Expr
    func E_desc(f string) *Expr
    func E_eq(f string, d interface{}) *Expr
    func E_field(f string) *Expr
    func E_ge(f string, d interface{}) *Expr
    func E_gt(f string, d interface{}) *Expr
    func E_in(f string, d interface{}) *Expr
    func E_is_null(f string) *Expr
    func E_le(f string, d interface{}) *Expr
    func E_like(f string, d string) *Expr
    func E_literal(d string, args ...interface{}) *Expr
    func E_lt(f string, d interface{}) *Expr
    func E_ne(f string, d interface{}) *Expr
    func E_not(d *Expr) *Expr
    func E_not_in(f string, d interface{}) *Expr
    func E_not_like(f string, d string) *Expr
    func E_not_null(f string) *Expr
    func E_or(d ...*Expr) *Expr
    func E_table(t string) *Expr
    func E_values(d ...interface{}) *Expr

// God是整个ORM的调度对象
type God struct{ ... }
    func NewGod(factory func() Model, node string, table Table) *God
    func NewGod_master_slave(factory func() Model, master_node string, slave_node string, table Table) *God
    func (g *God) Count() (int, error)
    func (g *God) Delete(args ...interface{}) (int, error)
    func (g *God) Exec(query string, args ...interface{}) (sql.Result, error)
    func (g *God) Groupby(args ...interface{}) *Statement
    func (g *God) Having(args ...*Expr) *Statement
    func (g *God) Insert(args ...interface{}) (int, error)
    func (g *God) Limit(args ...int) *Statement
    func (g *God) Load(pkv interface{}) (Model, error)
    func (g *God) Loads(pkvs ...interface{}) ([]Model, error)
    func (g *God) Multi(args ...interface{}) ([]Model, error)
    func (g *God) NewModel() Model
    func (g *God) NewStatement() *Statement
    func (g *God) On(args ...*Expr) *Statement
    func (g *God) Ondup(args ...*Expr) *Statement
    func (g *God) One(args ...interface{}) (Model, error)
    func (g *God) Orderby(args ...interface{}) *Statement
    func (g *God) Query(query string, args ...interface{}) (*sql.Rows, error)
    func (g *God) Sharding(datas ...interface{}) []*Statement
    func (g *God) Tpl(tpl string, fields ...interface{}) // 定义查询模板
    func (g *God) Update(args ...interface{}) (int, error)

// M和Model共同组成Model规范. 主要用来提供扩展数据的存储和方法
type Model interface{ ... }
type M struct{ ... }
    func (m *M) InitExtra()
    func (m *M) Extra(k string, v interface{})
    func (m *M) Extras(datas map[string]interface{})
    func (m *M) Exists(k string) bool
    func (m *M) Bool(k string) bool
    func (m *M) Float32(k string) float32
    func (m *M) Float64(k string) float64
    func (m *M) Int(k string) int
    func (m *M) Int16(k string) int16
    func (m *M) Int32(k string) int32
    func (m *M) Int64(k string) int64
    func (m *M) Int8(k string) int8
    func (m *M) Interface(k string) interface{}
    func (m *M) String(k string) string
    func (m *M) Uint(k string) uint
    func (m *M) Uint16(k string) uint16
    func (m *M) Uint32(k string) uint32
    func (m *M) Uint64(k string) uint64
    func (m *M) Uint8(k string) uint8

// Ref_xxx主要用于创建God对象时, Cache一些Model类型的反射信息
type Ref_Field struct{ ... }
    func NewRefField(model *Ref_Model, idx int) *Ref_Field
type Ref_Model struct{ ... }
    func NewRefModel(m interface{}) *Ref_Model
type Ref_Tag struct{ ... }

// Table和ShardingTable用于支持分表
type ShardingTable interface{ ... }
type Table interface{ ... }
    func Table_mod_int(format string, sharding_count int) Table
    func Table_normal(name string) Table
func (g *God) Sharding(datas ...interface{}) []*Statement
```

## 奇技淫巧

### 对象模板

```
GodOf_User.Tpl("simple", "id", "name", "passwd")
GodOf_User.On(E_eq("id", 1)).One("simple")
```

### 分表:多表查询

*核心思想*: 将分表规则应用到数据中, 产生独立的语句和分片数据对, 然后交由应用处理

```
var GodOf_User = NewGod(F_User, "数据库节点名", Table_mod_int("kk_user_%d", 256))

stmts := GodOf_User.Sharding([]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, []{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"})
for _, stmt := range stmts {
    stmt.On(
        E_in("id", stmt.ShardingData[0]),
        E_in("name", stmt.ShardingData[1]),
    ).Multi("simple")
}
```

### 分表:分表规则扩展

*核心思想*: 策略自身与策略的运行时, 仅通过Name(cols ...interface) string接口耦合

  * 每次对一行数据, 代入sharding相关列数据, 求表名

  * 多行数据的表名拆分, 由框架内部完成(上面的God.Sharding)

*扩展点*:

  1. 定义自己的分表策略对象, 字段根据自己策略需要的数据项来定

  2. 定义自己的策略工厂函数. 比如Table_mod_int

  3. 实现接口orm.Table { Name(cols ...interface{}) string }

```
type table_mod_int struct {
	Format        string
	ShardingCount int
}

func Table_mod_int(format string, sharding_count int) Table {
	if sharding_count < 1 {
		sharding_count = 1
	}
	return &table_mod_int{Format: format, ShardingCount: sharding_count}
}

func (t *table_mod_int) Name(cols ...interface{}) string {
	if len(cols) < 1 {
		return t.Format
	} else {
		switch v := cols[0].(type) {
		case int:
			return fmt.Sprintf(t.Format, int(v)%t.ShardingCount)
		case int64:
			return fmt.Sprintf(t.Format, int(v)%t.ShardingCount)
		case int32:
			return fmt.Sprintf(t.Format, int(v)%t.ShardingCount)
		case int16:
			return fmt.Sprintf(t.Format, int(v)%t.ShardingCount)
		case int8:
			return fmt.Sprintf(t.Format, int(v)%t.ShardingCount)
		case uint:
			return fmt.Sprintf(t.Format, int(v)%t.ShardingCount)
		case uint64:
			return fmt.Sprintf(t.Format, int(v)%t.ShardingCount)
		case uint32:
			return fmt.Sprintf(t.Format, int(v)%t.ShardingCount)
		case uint16:
			return fmt.Sprintf(t.Format, int(v)%t.ShardingCount)
		case uint8:
			return fmt.Sprintf(t.Format, int(v)%t.ShardingCount)
		default:
			return t.Format
		}
	}
}

func (t *table_mod_int) Names(cols ...interface{}) map[string][][]interface{} {
	return Names(t, cols...)
}
```