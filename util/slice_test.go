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

func TestCase_Slice_uniq_string(t *testing.T) {
	s1 := []string{"16628163", "28124695", "14556069", "116074", "18555476", "17383437", "11848177", "28123746", "25312803", "25429885", "28130783", "25673306", "14903406", "14556069", "116074", "18555476"}
	s2 := []string{"116074", "11848177", "14556069", "14903406", "16628163", "17383437", "18555476", "25312803", "25429885", "25673306", "28123746", "28124695", "28130783"}
	s3 := Slice_uniq_string(s1)
	if len(s2) != len(s3) || s2[3] != s3[3] || s2[5] != s3[5] {
		t.Logf("%#v", Slice_uniq_string(s1))
		t.Fail()
	}
}

func TestCase_Slice_string_remove(t *testing.T) {
	s := []string{"1", "22", "333", "4444", "55555"}
	sr := Slice_string_remove(s, "22")
	if len(sr) == len(s) || sr[0] != "1" || sr[1] != "333" || sr[2] != "4444" || sr[3] != "55555" {
		t.Fatal("slice_int_to_string failed")
	}
	s = []string{"1", "22", "333", "4444", "55555"}
	sr = Slice_string_remove(s, "222")
	if len(sr) != len(s) || sr[0] != "1" || sr[1] != "22" || sr[2] != "333" || sr[3] != "4444" || sr[4] != "55555" {
		t.Fatal("slice_int_to_string failed")
	}
	s = []string{"1", "22", "333", "4444", "55555"}
	sr = Slice_string_remove(s, "22", "333")
	if len(sr) == len(s) || sr[0] != "1" || sr[1] != "4444" || sr[2] != "55555" {
		t.Fatal("slice_int_to_string failed")
	}
}

func TestCase_Slice_string_remove_n(t *testing.T) {
	s := []string{"1", "22", "333", "4444", "55555"}
	sr := Slice_string_remove_n(s, 1, "22")
	if len(sr) == len(s) || sr[0] != "1" || sr[1] != "333" || sr[2] != "4444" || sr[3] != "55555" {
		t.Fatal("slice_int_to_string failed")
	}
	s = []string{"1", "22", "333", "4444", "55555"}
	sr = Slice_string_remove_n(s, 2, "22")
	if len(sr) == len(s) || sr[0] != "1" || sr[1] != "333" || sr[2] != "4444" || sr[3] != "55555" {
		t.Fatal("slice_int_to_string failed")
	}
	s = []string{"1", "22", "333", "4444", "55555"}
	sr = Slice_string_remove_n(s, 0, "22")
	if len(sr) != len(s) || sr[0] != "1" || sr[1] != "22" || sr[2] != "333" || sr[3] != "4444" || sr[4] != "55555" {
		t.Fatal("slice_int_to_string failed")
	}
	s = []string{"1", "22", "333", "4444", "55555"}
	sr = Slice_string_remove_n(s, 1, "55555")
	if len(sr) != 4 || sr[0] != "1" || sr[1] != "22" || sr[2] != "333" || sr[3] != "4444" {
		t.Fatal("slice_int_to_string failed")
	}
	s = []string{"1", "22", "333", "4444", "55555", "22"}
	sr = Slice_string_remove_n(s, 1, "22")
	if len(sr) == len(s) || sr[0] != "1" || sr[1] != "333" || sr[2] != "4444" || sr[3] != "55555" || sr[4] != "22" {
		t.Fatal("slice_int_to_string failed")
	}
	s = []string{"1", "22", "333", "4444", "55555", "22"}
	sr = Slice_string_remove_n(s, 2, "22")
	if len(sr) != 4 || sr[0] != "1" || sr[1] != "333" || sr[2] != "4444" || sr[3] != "55555" {
		t.Fatal("slice_int_to_string failed")
	}
	s = []string{"1", "22", "333", "4444", "55555", "22"}
	sr = Slice_string_remove_n(s, 2, "222")
	if len(sr) != len(s) || sr[0] != "1" || sr[1] != "22" || sr[2] != "333" || sr[3] != "4444" || sr[4] != "55555" || sr[5] != "22" {
		t.Fatal("slice_int_to_string failed")
	}
}
