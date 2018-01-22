package util

import (
	"reflect"
	"testing"
)

type M_Indirect_User struct {
	Name string
	Age  interface{}
}

func Test_Indirect(t *testing.T) {
	var (
		v1 *****M_Indirect_User
		v2 M_Indirect_User = M_Indirect_User{Name: "HELLO"}
		rv reflect.Value
	)
	rv = Indirect(&v1)
	if (*****v1).Name != "" || v2.Name != "HELLO" || rv.Interface().(M_Indirect_User).Name != "" {
		t.Logf("Before change value:")
		t.Logf("*****v1: %#v", *****v1)
		t.Logf("v2: %#v", v2)
		t.Logf("rv: %#v", rv)
	}

	rv.Set(reflect.ValueOf(v2))
	if (*****v1).Name != "HELLO" || v2.Name != "HELLO" || rv.Interface().(M_Indirect_User).Name != "HELLO" {
		t.Logf("After Set v2 into rv:")
		t.Logf("*****v1: %#v", *****v1)
		t.Logf("v2: %#v", v2)
		t.Logf("rv: %#v", rv)
	}

	v2.Name = "Hello World"
	if (*****v1).Name != "HELLO" || v2.Name != "Hello World" || rv.Interface().(M_Indirect_User).Name != "HELLO" {
		t.Logf("Modify v2.Name:")
		t.Logf("*****v1: %#v", *****v1)
		t.Logf("v2: %#v", v2)
		t.Logf("rv: %#v", rv)
		t.Fail()
	}
}

func Test_Indirect_stopAt(t *testing.T) {
	var (
		v1 *****M_Indirect_User
		v2 M_Indirect_User = M_Indirect_User{Name: "HELLO"}
		rv reflect.Value
	)
	rv = Indirect_stopAt(&v1, func(rv reflect.Value) bool {
		return rv.Kind() == reflect.Ptr && rv.Elem().Kind() == reflect.Struct
	})
	if (*****v1).Name != "" || v2.Name != "HELLO" || rv.Interface().(*M_Indirect_User).Name != "" {
		t.Logf("Before change value:")
		t.Logf("*****v1: %#v", *****v1)
		t.Logf("v2: %#v", v2)
		t.Logf("rv: %#v", rv)
		t.Fail()
	}

	rv.Set(reflect.ValueOf(&v2))
	if (*****v1).Name != "HELLO" || v2.Name != "HELLO" || rv.Interface().(*M_Indirect_User).Name != "HELLO" {
		t.Logf("After Set v2 into rv:")
		t.Logf("*****v1: %#v", *****v1)
		t.Logf("v2: %#v", v2)
		t.Logf("rv: %#v", rv)
		t.Fail()
	}

	v2.Name = "Hello World"
	if (*****v1).Name != "Hello World" || v2.Name != "Hello World" || rv.Interface().(*M_Indirect_User).Name != "Hello World" {
		t.Logf("Modify v2.Name:")
		t.Logf("*****v1: %#v", *****v1)
		t.Logf("v2: %#v", v2)
		t.Logf("rv: %#v", rv)
		t.Fail()
	}
}

func Test_Indirect_zero(t *testing.T) {
	var (
		v1       M_Indirect_User
		rv1, rv2 reflect.Value
		iage     interface{}
		age      ******int
	)
	iage = age
	v1.Age = &iage
	rv1 = Indirect_zero(&v1.Age, false)
	rv2 = Indirect_zero(&v1.Age, true)
	if _, ok := rv1.Interface().(******int); ok {
		if _, ok := rv2.Interface().(*interface{}); ok {
			return
		}
	}

	t.Logf("type of IndireIndirect_zero((&v1.Age, false): %T", rv1.Interface())
	t.Logf("type of IndireIndirect_zero((&v1.Age, true): %T", rv2.Interface())
	t.Fail()
}
