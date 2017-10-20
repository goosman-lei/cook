package model

import (
	"crypto/md5"
	"fmt"
	"strconv"
)

const (
	SHARD_TYPE_NONE uint8 = iota
	SHARD_TYPE_MOD_INT
	SHARD_TYPE_MOD_STRING
)

func (g *God) shard(input_args ...interface{}) string {
	switch g.Opts.ShardType {
	case SHARD_TYPE_MOD_STRING:
		return g.shard_mod_string(input_args...)
	case SHARD_TYPE_MOD_INT:
		return g.shard_mod_int(input_args...)
	case SHARD_TYPE_NONE:
		return g.shard_none(input_args...)
	default:
		return g.shard_none(input_args...)
	}
}

func (g *God) shard_none(input_args ...interface{}) string {
	return g.Opts.TableFmt
}

func (g *God) shard_mod_int(input_args ...interface{}) string {
	if len(input_args) > 0 {
		switch v := input_args[0].(type) {
		case int:
			return fmt.Sprintf(g.Opts.TableFmt, v%g.Opts.ShardCnt)
		case int8:
			return fmt.Sprintf(g.Opts.TableFmt, v%int8(g.Opts.ShardCnt))
		case int16:
			return fmt.Sprintf(g.Opts.TableFmt, v%int16(g.Opts.ShardCnt))
		case int32:
			return fmt.Sprintf(g.Opts.TableFmt, v%int32(g.Opts.ShardCnt))
		case int64:
			return fmt.Sprintf(g.Opts.TableFmt, v%int64(g.Opts.ShardCnt))
		case uint:
			return fmt.Sprintf(g.Opts.TableFmt, v%uint(g.Opts.ShardCnt))
		case uint8:
			return fmt.Sprintf(g.Opts.TableFmt, v%uint8(g.Opts.ShardCnt))
		case uint16:
			return fmt.Sprintf(g.Opts.TableFmt, v%uint16(g.Opts.ShardCnt))
		case uint32:
			return fmt.Sprintf(g.Opts.TableFmt, v%uint32(g.Opts.ShardCnt))
		case uint64:
			return fmt.Sprintf(g.Opts.TableFmt, v%uint64(g.Opts.ShardCnt))
		}
	}
	return g.shard_none()
}

func (g *God) shard_mod_string(input_args ...interface{}) string {
	i_cols := make([]interface{}, len(input_args))
	for i, col := range input_args {
		v := int64(0)
		switch g.Opts.HashAlgo {
		default:
			v, _ = strconv.ParseInt(fmt.Sprintf("%x\n", md5.Sum([]byte(col.(string))))[0:8], 16, 64)
		}
		i_cols[i] = int(v) % g.Opts.ShardCnt
	}
	return fmt.Sprintf(g.Opts.TableFmt, i_cols...)
}
