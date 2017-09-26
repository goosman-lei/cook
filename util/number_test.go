package util

import (
	"sort"
	"testing"
)

func TestCase_Uniq_int(t *testing.T) {
	s := Uniq_int(sort.IntSlice{1, 2, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4, 5, 5})
	if len(s) != 5 || s[0] != 1 || s[2] != 3 || s[4] != 5 {
		t.Fatal("uniq_int failed")
	}
}
