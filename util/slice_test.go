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

type M1 struct {
	Id   int
	Name string
}
type M2 struct {
	Id   int
	Name string
}

func TestCase_Slice_pick(t *testing.T) {
	s1 := []M1{M1{1, "Jack"}, M1{2, "Tom"}, M1{3, "Green"}, M1{4, "Tim"}, M1{5, "Tony"}}
	s2 := &[]M2{M2{1, "Jack"}, M2{2, "Tom"}, M2{3, "Green"}, M2{4, "Tim"}, M2{5, "Tony"}}
	s3 := []*M2{&M2{1, "Jack"}, &M2{2, "Tom"}, &M2{3, "Green"}, &M2{4, "Tim"}, &M2{5, "Tony"}}
	s4 := &[]*M2{&M2{1, "Jack"}, &M2{2, "Tom"}, &M2{3, "Green"}, &M2{4, "Tim"}, &M2{5, "Tony"}}

	r_ints := Slice_pick_int(s1, "Id")
	if len(r_ints) != 5 || r_ints[2] != 3 || r_ints[4] != 5 {
		t.Logf("%v", r_ints)
		t.Fail()
	}

	r_ints = Slice_pick_int(s2, "Id")
	if len(r_ints) != 5 || r_ints[2] != 3 || r_ints[4] != 5 {
		t.Logf("%v", r_ints)
		t.Fail()
	}

	r_ints = Slice_pick_int(s3, "Id")
	if len(r_ints) != 5 || r_ints[2] != 3 || r_ints[4] != 5 {
		t.Logf("%v", r_ints)
		t.Fail()
	}

	r_ints = Slice_pick_int(s4, "Id")
	if len(r_ints) != 5 || r_ints[2] != 3 || r_ints[4] != 5 {
		t.Logf("%v", r_ints)
		t.Fail()
	}

	r_strings := Slice_pick_string(s1, "Name")
	if len(r_strings) != 5 || r_strings[2] != "Green" || r_strings[4] != "Tony" {
		t.Logf("%v", r_strings)
		t.Fail()
	}
}
