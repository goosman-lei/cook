package option

import (
	"testing"
)

type TmpOptions struct {
	Name string
	Age  int
	A    string
	B    string
	C    string
	D    string
	E    string
	F    string
	G    string
	H    string
	I    string
	J    string
}

func Test_Opt_string(t *testing.T) {
	var (
		ok   bool
		err  error
		oerr OptionError
		opts *TmpOptions = &TmpOptions{}
	)
	if err = Apply(opts, Opt("Name", "Jack"), Opt("Age", 10)); err != nil {
		t.Logf("unexpect error: %s\n", err)
		t.Fail()
	}
	if opts.Name != "Jack" || opts.Age != 10 {
		t.Logf("unexpect result option: %#v", opts)
		t.Fail()
	}
	if err = Apply(opts, Opt("Age", 20)); err != nil {
		t.Logf("unexpect error: %s\n", err)
		t.Fail()
	}
	if opts.Name != "Jack" || opts.Age != 20 {
		t.Logf("unexpect result option: %#v", opts)
		t.Fail()
	}
	if err = Apply(opts, Opt("Age", 20)); err != nil {
		t.Logf("unexpect error: %s\n", err)
		t.Fail()
	}

	if err = Apply(opts, Opt("Aget", 30)); err == nil {
		t.Logf("expect error, but have no\n")
		t.Fail()
	}
	if oerr, ok = err.(OptionError); !ok {
		t.Logf("expect OptionError, but: %s\n", err)
		t.Fail()
	}
	if oerr.ErrType != Err_Field_Invalid {
		t.Logf("expect field invalid error, but: %s\n", err)
		t.Fail()
	}
}

func Benchmark_Opt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Opt("Name", "Jack")
	}
}

func Benchmark_Apply(b *testing.B) {
	t := &TmpOptions{}
	o := Opt("Name", "Jack")
	for i := 0; i < b.N; i++ {
		Apply(t, o)
	}
}

// apply multi's real-time used is linear growth
func Benchmark_Apply_Multi(b *testing.B) {
	t := &TmpOptions{}
	o := []Option{
		Opt("A", "Jack"),
		Opt("B", "Jack"),
		Opt("C", "Jack"),
		Opt("D", "Jack"),
		Opt("E", "Jack"),
		Opt("F", "Jack"),
		Opt("G", "Jack"),
		Opt("H", "Jack"),
		Opt("I", "Jack"),
		Opt("J", "Jack"),
	}
	for i := 0; i < b.N; i++ {
		Apply(t, o...)
	}
}
