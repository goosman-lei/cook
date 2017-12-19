package connector

import (
	"github.com/garyburd/redigo/redis"
	"net"
	"reflect"
	"testing"
	"time"
	"unsafe"
)

var (
	testAddr_1 string = "10.10.200.10:6850"
	testAddr_2 string = "10.10.200.10:6857"
	skipRedis  bool   = false
)

func init() {
	if _, err := net.DialTimeout("tcp", testAddr_1, 1000*time.Millisecond); err != nil {
		skipRedis = true
		return
	}
	if _, err := net.DialTimeout("tcp", testAddr_2, 1000*time.Millisecond); err != nil {
		skipRedis = true
		return
	}

	configs := map[string]RedisConf{
		"cluster1": RedisConf{
			Addrs:          []string{testAddr_1},
			TestInterval:   time.Minute,
			MaxActive:      16,
			MaxIdle:        16,
			IdleTimeout:    5 * time.Minute,
			ConnectTimeout: 1000 * time.Millisecond,
			ReadTimeout:    1000 * time.Millisecond,
			WriteTimeout:   1000 * time.Millisecond,
		},
		"cluster2": RedisConf{
			Addrs:          []string{testAddr_2},
			TestInterval:   time.Minute,
			MaxActive:      16,
			MaxIdle:        16,
			IdleTimeout:    5 * time.Minute,
			ConnectTimeout: 1000 * time.Millisecond,
			ReadTimeout:    1000 * time.Millisecond,
			WriteTimeout:   1000 * time.Millisecond,
		},
	}

	SetupRedis(configs)
}

func Test_GetRedis(t *testing.T) {
	if skipRedis {
		t.Skipf("test redis server is not reachable")
	}
	c1 := MustGetRedis("cluster1")
	c2 := MustGetRedis("cluster2")

	r_c1_conn_ptr := reflect.ValueOf(c1.Conn).Elem().FieldByName("c").Elem().Elem().FieldByName("conn").Elem().Pointer()
	r_c1_conn := (*net.TCPConn)(unsafe.Pointer(r_c1_conn_ptr))
	if r_c1_conn.RemoteAddr().String() != testAddr_1 {
		t.Logf("cluster1 address wrong: want %s, actually %s", testAddr_1, r_c1_conn.RemoteAddr().String())
		t.Fail()
	}

	r_c2_conn_ptr := reflect.ValueOf(c2.Conn).Elem().FieldByName("c").Elem().Elem().FieldByName("conn").Elem().Pointer()
	r_c2_conn := (*net.TCPConn)(unsafe.Pointer(r_c2_conn_ptr))
	if r_c2_conn.RemoteAddr().String() != testAddr_2 {
		t.Logf("cluster2 address wrong: want %s, actually %s", testAddr_2, r_c2_conn.RemoteAddr().String())
		t.Fail()
	}

}

func TestRedis(t *testing.T) {
	if skipRedis {
		t.Skipf("test redis server is not reachable")
	}

	var (
		err     error
		strVal  string
		boolVal bool
		strVals []string
		intVal  int
	)

	strVal, err = MustGetRedis("cluster1").DoString("SET", "overmind-ng-test-key:a", "HELLO-WOR")
	if err != nil {
		t.Logf("(SET overmind-ng-test-key:a 1) unexptected error: %s", err)
		t.Fail()
	}
	if strVal != "OK" {
		t.Logf("(SET overmind-ng-test-key:a 1) unexptected return: %t", strVal)
		t.Fail()
	}

	strVal, err = MustGetRedis("cluster1").DoString("SET", "overmind-ng-test-key:a", "HELLO-WORLD", "NX")
	if err == nil {
		t.Logf("(SET overmind-ng-test-key:a 1 NX) expect error, but nil")
		t.Fail()
	}

	strVal, err = MustGetRedis("cluster1").DoString("GET", "overmind-ng-test-key:a")
	if err != nil {
		t.Logf("(GET overmind-ng-test-key:a) unexptected error: %s", err)
		t.Fail()
	}
	if strVal != "HELLO-WOR" {
		t.Logf("(GET overmind-ng-test-key:a) unexptected return: %s", strVal)
		t.Fail()
	}

	boolVal, err = MustGetRedis("cluster1").DoBool("EXISTS", "overmind-ng-test-key:a")
	if err != nil {
		t.Logf("(EXISTS overmind-ng-test-key:a) unexptected error: %s", err)
		t.Fail()
	}
	if !boolVal {
		t.Logf("(EXISTS overmind-ng-test-key:a) unexptected return: %t", boolVal)
		t.Fail()
	}

	boolVal, err = MustGetRedis("cluster1").DoBool("EXISTS", "overmind-ng-test-key:b")
	if err != nil {
		t.Logf("(EXISTS overmind-ng-test-key:b) unexptected error: %s", err)
		t.Fail()
	}
	if boolVal {
		t.Logf("(EXISTS overmind-ng-test-key:b) unexptected return: %t", boolVal)
		t.Fail()
	}

	strVals, err = MustGetRedis("cluster1").DoStrings("MGET", "overmind-ng-test-key:a", "overmind-ng-test-key:b")
	if err != nil {
		t.Logf("(MGET overmind-ng-test-key:a overmind-ng-test-key:b) unexptected error: %s", err)
		t.Fail()
	}
	if len(strVals) != 2 || strVals[0] != "HELLO-WOR" || strVals[1] != "" {
		t.Logf("(MGET overmind-ng-test-key:a overmind-ng-test-key:b) unexptected return: %q", strVals)
		t.Fail()
	}

	intVal, err = redis.Int(MustGetRedis("cluster1").Do("DEL", "overmind-ng-test-key:a"))
	if err != nil {
		t.Logf("(DEL overmind-ng-test-key:a) unexptected error: %s", err)
		t.Fail()
	}
	if intVal != 1 {
		t.Logf("(DEL overmind-ng-test-key:a) unexptected return: %q", intVal)
		t.Fail()
	}

}
