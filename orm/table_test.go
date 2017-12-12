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

type M_User_sharding_test struct {
	*M
}

func F_User_sharding_test() Model { return &M_User_sharding_test{} }

func Test_Table_Sharding(t *testing.T) {
	t.Skip()
	var God = NewGod(F_User_sharding_test, "not-exists-node", Table_mod_int("kk_user_%d", 6))

	stmts := God.Shardings([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	stmts[0].On(E_in("id", stmts[0].ShardingData[0])).One()
	if stmts[0].SQL != "SELECT * FROM kk_user_3 WHERE id IN(?, ?, ?) LIMIT ?" || len(stmts[0].Args) != 4 || stmts[0].Args[0] != 3 {
		t.Logf("error at sharding:\n\tSQL: %s\n\tArgs: %#v\n", stmts[0].SQL, stmts[0].Args)
		t.Fail()
	}
}
