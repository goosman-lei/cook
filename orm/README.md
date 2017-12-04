# cook/orm: 从互联网应用架构角度出发的ORM框架

## 特性

* 主从分离

* 水平分表

* 对象映射

* 结构化查询

* 调试接口

* 不支持"对象关系映射(描述对象关系转化JOIN)", 原因如下

    * 移动互联网依赖, 大大提升了server端需要处理的业务量

    * 大访问量下, 我们通常会做缓存/分表等架构方案, 这种情况下, 对象关系映射反而成为问题

## API

### Register model's God into Heaven

type GodOptions struct {
    Model_type reflect.Type

    Node string
    
}


func NewGod(model interface{}, table Table)

### In God View

```
func (g *God) Load(model interface{}, pk interface{}) error

func (g *God) Loads(models interface{}, pk []interface{}) error

func (g *God) Tpl(tpl string) *Query

func (g *God) Cols(cols ...interface{}) *Query

func (g *God) Where(exprs ...*Expr) *Query

func (g *God) GroupBy(exprs ...*Expr) *Query

func (g *God) Having(exprs ...*Expr) *Query

func (g *God) OrderBy(exprs ...*Expr) *Query

func (g *God) Offset(offset int) *Query

func (g *God) Limit(limit int) *Query

func (g *God) Set(exprs ...*Expr) *Query

func (g *God) OnDup(exprs ...*Expr) *Query

func (g *God) Values(exprs ...*Expr) *Query

func (g *God) Count(cols ...interface{}) (int, error)

func (g *God) One(model interface{}) error

func (g *God) Multi(models interface{}) error

func (g *God) Delete() (int, error)

func (g *God) Update() (int, error)

func (g *God) Insert() (int, int, error)
```