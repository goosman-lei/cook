package util

import (
	"fmt"
	"testing"
	"time"
)

func Test_Concurrent(t *testing.T) {
	m := NewCMap()
	done := make(chan bool)

	go func() {
		for i := 0; i < 50000; i++ {
			m.Set(fmt.Sprintf("key-%d", i), i)
		}
		done <- true
	}()

	go func() {
		time.Sleep(1e6)
		for i := 0; i < 50000; i++ {
			m.Get(fmt.Sprintf("key-%d", i))
		}
		done <- true
	}()
	<-done
	<-done
}

func Test_Iterate(t *testing.T) {
	m := NewCMap()
	m.Set("a", 0)
	m.Set("b", 1)
	m.Set("c", 2)
	m.Set("d", 3)
	m.Set("e", 4)
	m.Set("f", 5)
	m.Set("g", 6)

	s := make([]string, 7)
	m.Iterate(func(k string, v interface{}) {
		s[v.(int)] = k
	})
	if s[0] != "a" || s[1] != "b" || s[2] != "c" || s[3] != "d" || s[4] != "e" || s[5] != "f" || s[6] != "g" {
		t.Logf("iterate result unexpected: %q", s)
		t.Fail()
	}

	m.Apply(func(k string, v interface{}) interface{} {
		return v.(int) * 2
	})
	if m.MustGet("a") != 0 || m.MustGet("b") != 2 || m.MustGet("c") != 4 || m.MustGet("d") != 6 || m.MustGet("e") != 8 || m.MustGet("f") != 10 || m.MustGet("g") != 12 {
		t.Logf("apply result unexpected: %q", m)
		t.Fail()
	}
}

func BenchmarkCase_Set(b *testing.B) {
	m := NewCMap()
	for i := 0; i < 10000; i++ {
		m.Set(fmt.Sprintf("KEY%d", i), i)
	}
	for i := 0; i < b.N; i++ {
		m.Set("KEY", i)
	}
}

func BenchmarkCase_SetAndErase(b *testing.B) {
	m := NewCMap()
	for i := 0; i < b.N; i++ {
		m.Erase("KEY")
	}
}

func BenchmarkCase_Get(b *testing.B) {
	m := NewCMap()
	m.Set("KEY", "VAL")
	for i := 0; i < b.N; i++ {
		m.Get("KEY")
	}
}
