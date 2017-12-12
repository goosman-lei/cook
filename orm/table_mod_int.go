package orm

import (
	"fmt"
	cook_util "gitlab.niceprivate.com/golang/cook/util"
)

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
		return fmt.Sprintf(t.Format, cook_util.As_int(cols[0])%t.ShardingCount)
	}
}

func (t *table_mod_int) Names(cols ...interface{}) map[string][][]interface{} {
	return Names(t, cols...)
}
