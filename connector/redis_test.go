package connector

import (
	"github.com/garyburd/redigo/redis"
	"net"
	"testing"
	"time"
)

var skipRedis bool = false

func init() {
	testAddr := "10.10.200.10:6850"
	if _, err := net.DialTimeout("tcp", testAddr, 1000*time.Millisecond); err != nil {
		skipRedis = true
		return
	}

	configs := map[string]RedisConf{
		"cluster1": RedisConf{
			Addrs:          []string{testAddr},
			TestInterval:   time.Minute,
			MaxActive:      16,
			MaxIdle:        16,
			IdleTimeout:    5 * time.Minute,
			ConnectTimeout: 1000 * time.Millisecond,
			ReadTimeout:    1000 * time.Millisecond,
			WriteTimeout:   1000 * time.Millisecond,
		},
		"cluster2": RedisConf{
			Addrs:          []string{testAddr},
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
