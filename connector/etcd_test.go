package connector

import (
	"sort"
	"testing"
	"time"
)

func init() {
	SetupEtcd(EtcdConf{
		Addrs:   []string{"http://10.10.10.29:2379", "http://10.10.10.30:2379", "http://10.10.10.70:2379"},
		Timeout: time.Second,
	})
}

func Test_etcd_kv(t *testing.T) {
	if err := Etcd.Remove("/ut-etcd-prefix/direct-key"); err != nil {
		t.Logf("remove failed: %s", err)
		t.Fail()
	}

	if err := Etcd.Set("/ut-etcd-prefix/direct-key", "10.10.200.12:2020"); err != nil {
		t.Logf("set failed: %s", err)
		t.Fail()
	}

	if v, err := Etcd.Get("/ut-etcd-prefix/direct-key"); err != nil {
		t.Logf("get failed: %s", err)
		t.Fail()
	} else if v != "10.10.200.12:2020" {
		t.Logf("get wrong value: %s, want: 10.10.200.12:2020", v)
		t.Fail()
	}

	if err := Etcd.Remove("/ut-etcd-prefix/direct-key"); err != nil {
		t.Logf("remove failed: %s", err)
		t.Fail()
	}
}

func Test_etcd_dir(t *testing.T) {
	if err := Etcd.Remove("/ut-etcd-prefix"); err != nil {
		t.Logf("remove failed: %s", err)
		t.Fail()
	}

	if err := Etcd.Mkdir("/ut-etcd-prefix/a/b/c"); err != nil {
		t.Logf("mkdir failed: %s", err)
		t.Fail()
	}

	if err := Etcd.Mkdir("/ut-etcd-prefix/a/b1/c1"); err != nil {
		t.Logf("mkdir failed: %s", err)
		t.Fail()
	}

	if err := Etcd.Mkdir("/ut-etcd-prefix/a2/b/c"); err != nil {
		t.Logf("mkdir failed: %s", err)
		t.Fail()
	}

	if err := Etcd.Set("/ut-etcd-prefix/a3/value", "10.10.200.12:2020"); err != nil {
		t.Logf("set failed: %s", err)
		t.Fail()
	}

	if v, err := Etcd.Dir("/ut-etcd-prefix"); err != nil {
		t.Logf("get failed: %s", err)
		t.Fail()
	} else {
		sv := sort.StringSlice(v)
		sv.Sort()
		if sv[0] != "/ut-etcd-prefix/a" || sv[1] != "/ut-etcd-prefix/a2" || sv[2] != "/ut-etcd-prefix/a3" {
			t.Logf("wrong dir values: %#v", v)
			t.Fail()
		}
	}

	if v, err := Etcd.IsDir("/ut-etcd-prefix"); err != nil {
		t.Logf("get failed: %s", err)
		t.Fail()
	} else if !v {
		t.Logf("is_dir check error")
		t.Fail()
	}

	if v, err := Etcd.IsDir("/ut-etcd-prefix/a3/value"); err != nil {
		t.Logf("get failed: %s", err)
		t.Fail()
	} else if v {
		t.Logf("is_dir check error")
		t.Fail()
	}

	if err := Etcd.Remove("/ut-etcd-prefix"); err != nil {
		t.Logf("remove failed: %s", err)
		t.Fail()
	}
}
