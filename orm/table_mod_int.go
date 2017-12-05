package orm

import (
	"fmt"
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
