package orm

import (
	"fmt"
	cook_util "gitlab.niceprivate.com/golang/cook/util"
	"strconv"
)

type table_hash_p8_mod_int struct {
	Format        string
	ShardingCount int
}

func Table_hash_p8_mod_int(format string, sharding_count int) Table {
	if sharding_count < 1 {
		sharding_count = 1
	}
	return &table_hash_p8_mod_int{Format: format, ShardingCount: sharding_count}
}

func (t *table_hash_p8_mod_int) Name(cols ...interface{}) string {
	if len(cols) < 1 {
		return t.Format
	} else if iv, err := strconv.ParseInt(cook_util.As_string(cols[0])[0:8], 16, 64); err != nil {
		return fmt.Sprintf(t.Format, 0)
	} else {
		return fmt.Sprintf(t.Format, int(iv)%t.ShardingCount)
	}
}

func (t *table_hash_p8_mod_int) Names(cols ...interface{}) map[string][][]interface{} {
	return Names(t, cols...)
}
