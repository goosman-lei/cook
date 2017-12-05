package orm

import (
	"fmt"
	"testing"
)

func Test_Table_normal(t *testing.T) {
	table := Table_normal("kk_user")

	if table.Name() != "kk_user" {
		t.Logf("unexcept table name: %s", table.Name())
		t.Fail()
	}

	if table.Name(1, 2, 3) != "kk_user" {
		t.Logf("unexcept table name: %s", table.Name())
		t.Fail()
	}
}

func Test_Table_mod_int(t *testing.T) {
	table := Table_mod_int("kk_user_%d", 128)

	if table.Name(5012470) != fmt.Sprintf("kk_user_%d", 5012470%128) {
		t.Logf("unexcept table name: %s", table.Name(5012470))
		t.Fail()
	}

	if table.Name(5012470, 1, 2, 3) != fmt.Sprintf("kk_user_%d", 5012470%128) {
		t.Logf("unexcept table name: %s", table.Name(5012470, 1, 2, 3))
		t.Fail()
	}

	shardingTable := table.(ShardingTable)

	tableNames := shardingTable.Names([]int{1, 2, 3, 4, 5, 129, 130})
	if len(tableNames) != 5 {
		t.Logf("unexcept shardingTable name: %s", tableNames)
		t.Fail()
	}
}
