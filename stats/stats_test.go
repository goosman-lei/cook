package stats

import (
	"net/http"
	_ "net/http/pprof"
	"testing"
	"time"
)

func init() {
	Init_stats()
}

func Test_Concurrence_incr_and_decr(t *testing.T) {
	done := make(chan bool)
	mapping = make(map[string]interface{})
	update()
	if Eye != "{}" {
		t.Logf("Eye init value invalid")
		t.Fail()
	}

	go func() {
		http.ListenAndServe(":6060", nil)
	}()

	nGoroutine := 100
	nLoop := 100

	for i := 0; i < nGoroutine; i++ {
		go func() {
			for j := 0; j < nLoop; j++ {
				Incr_int("testing.incr_and_decr", 1)
			}
			done <- true
		}()
	}
	for i := 0; i < nGoroutine; i++ {
		go func() {
			for j := 0; j < nLoop; j++ {
				Decr_int("testing.incr_and_decr", 1)
			}
			done <- true
		}()
	}

	for i := 0; i < 2*nGoroutine; i++ {
		<-done
	}

	time.Sleep(1e8)
	map_mutex.RLock()
	defer map_mutex.RUnlock()
	if mapping["testing"].(map[string]interface{})["incr_and_decr"].(int) != 0 {
		t.Fatalf("testing.incr_and_decr want 0 have %d", mapping["testing"].(map[string]interface{})["incr_and_decr"].(int))
	}
}

func Test_String(t *testing.T) {
	mapping = make(map[string]interface{})
	update()
	if Eye != "{}" {
		t.Logf("Eye init value invalid")
		t.Fail()
	}

	set("china.beijing.haozan.leiguoguo.name", "goosman-lei")
	update()
	if Eye != "{china: {beijing: {haozan: {leiguoguo: {name: \"goosman-lei\"}}}}}" {
		t.Fatalf("Eye value invalid: %s", Eye)
	}
}

func Test_Bool(t *testing.T) {
	mapping = make(map[string]interface{})
	update()
	if Eye != "{}" {
		t.Logf("Eye init value invalid")
		t.Fail()
	}

	set("china.beijing.haozan.leiguoguo.is_old", true)
	update()
	if Eye != "{china: {beijing: {haozan: {leiguoguo: {is_old: true}}}}}" {
		t.Fatalf("Eye value invalid: %s", Eye)
	}
}

func Test_Strings(t *testing.T) {
	mapping = make(map[string]interface{})
	update()
	if Eye != "{}" {
		t.Logf("Eye init value invalid")
		t.Fail()
	}

	set("china.beijing.haozan.leiguoguo.family", []string{"father", "mother", "brother", "sister", "myself", "wife", "son", "niece", "nephew"})
	update()
	if Eye != "{china: {beijing: {haozan: {leiguoguo: {family: [\"father\", \"mother\", \"brother\", \"sister\", \"myself\", \"wife\", \"son\", \"niece\", \"nephew\"]}}}}}" {
		t.Fatalf("Eye value invalid: %s", Eye)
	}
}

func Test_Ints(t *testing.T) {
	mapping = make(map[string]interface{})
	update()
	if Eye != "{}" {
		t.Logf("Eye init value invalid")
		t.Fail()
	}

	set("china.beijing.haozan.leiguoguo.nums", []int{1, 2, 3, 8, 9, 4, 5, 7, 6, 0})
	update()
	if Eye != "{china: {beijing: {haozan: {leiguoguo: {nums: [1, 2, 3, 8, 9, 4, 5, 7, 6, 0]}}}}}" {
		t.Fatalf("Eye value invalid: %s", Eye)
	}
}

func Test_Uint32s(t *testing.T) {
	mapping = make(map[string]interface{})
	update()
	if Eye != "{}" {
		t.Logf("Eye init value invalid")
		t.Fail()
	}

	set("china.beijing.haozan.leiguoguo.nums", []uint32{1, 2, 3, 8, 9, 4, 5, 7, 6, 0})
	update()
	if Eye != "{china: {beijing: {haozan: {leiguoguo: {nums: [1, 2, 3, 8, 9, 4, 5, 7, 6, 0]}}}}}" {
		t.Fatalf("Eye value invalid: %s", Eye)
	}
}

func Test_Int(t *testing.T) {
	mapping = make(map[string]interface{})
	update()
	if Eye != "{}" {
		t.Fatalf("Eye init value invalid")
	}

	incr("china.beijing.haozan.leiguoguo.age", 1)
	update()
	if Eye != "{china: {beijing: {haozan: {leiguoguo: {age: 1}}}}}" {
		t.Fatalf("Eye value invalid: %s", Eye)
	}
	incr("china.beijing.haozan.leiguoguo.age", 1)
	incr("china.beijing.haozan.leiguoguo.age", 1)
	incr("china.beijing.haozan.leiguoguo.age", 1)
	incr("china.beijing.haozan.leiguoguo.age", 1)
	update()
	if Eye != "{china: {beijing: {haozan: {leiguoguo: {age: 5}}}}}" {
		t.Fatalf("Eye value invalid: %s", Eye)
	}
	decr("china.beijing.haozan.leiguoguo.age", 1)
	decr("china.beijing.haozan.leiguoguo.age", 1)
	update()
	if Eye != "{china: {beijing: {haozan: {leiguoguo: {age: 3}}}}}" {
		t.Fatalf("Eye value invalid: %s", Eye)
	}
}

func Benchmark_set(b *testing.B) {
	refresh_interval = time.Hour * 24
	for i := 0; i < b.N; i++ {
		set("category.string_key", "goosman-lei")
	}
}

func Benchmark_incr(b *testing.B) {
	refresh_interval = time.Hour * 24
	for i := 0; i < b.N; i++ {
		incr("category.int_key", 1)
	}
}

func Benchmark_decr(b *testing.B) {
	refresh_interval = time.Hour * 24
	for i := 0; i < b.N; i++ {
		decr("category.int_key", 1)
	}
}

func Benchmark_Incr_Int(b *testing.B) {
	refresh_interval = time.Hour * 24
	for i := 0; i < b.N; i++ {
		Incr_int("category.int_incr_key", 1)
	}
}

func Benchmark_update_complex(b *testing.B) {
	refresh_interval = time.Hour * 24
	set("china.beijing.haozan.leiguoguo.name", "goosman-lei")
	set("china.beijing.haozan.leiguoguo.age", 32)
	set("english.nums", []string{"one", "two", "three", "four", "five", "nice", "eight", "seven", "six", "zero"})
	set("arabia.nums", []uint32{1, 2, 3, 8, 9, 4, 5, 7, 6, 0})
	set("china.beijing.haozan.leiguoguo.family", []string{"father", "mother", "brother", "sister", "myself", "wife", "son", "niece", "nephew"})
	for i := 0; i < b.N; i++ {
		update()
	}
}

func Benchmark_update_normal(b *testing.B) {
	refresh_interval = time.Hour * 24
	startupTime := time.Now().Truncate(time.Hour * 24 * 30)
	set("core.running_time", func() string {
		return time.Now().Sub(startupTime).String()
	})
	set("gate.total_client_num", 1000000)
	set("gate.online_client_num", 3728)
	set("gate.total_send_msg", 173720382)
	set("gate.fail_send_msg", 183)
	set("gate.total_recv_msg", 182390173818)
	set("gate.fail_recv_msg", 17328)
	set("proxy.total_client_num", 1000000)
	set("proxy.online_client_num", 3728)
	set("proxy.total_send_msg", 173720382)
	set("proxy.fail_send_msg", 183)
	set("proxy.total_recv_msg", 182390173818)
	set("proxy.fail_recv_msg", 17328)
	for i := 0; i < b.N; i++ {
		update()
	}
}
