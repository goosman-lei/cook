package util

import (
	"testing"
)

func TestCase_Slice_int_to_string(t *testing.T) {
	s := []int{1, 22, 333, 4444, 55555}
	sr := Slice_int_to_string(s)
	if len(sr) != len(s) || sr[0] != "1" || sr[2] != "333" || sr[4] != "55555" {
		t.Fatal("slice_int_to_string failed")
	}
}
