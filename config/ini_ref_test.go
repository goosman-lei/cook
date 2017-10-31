package config

import (
	"testing"
)

func Test_ini_ref_config(t *testing.T) {
	if err := init_ref_config("ut-conf"); err != nil {
		t.Fatalf("parse ref config failed: %s", err)
	}
	if ref_config["path.home"] != "/home/nice" || ref_config["path.var"] != "/home/nice/var" || ref_config["path.log"] != "/home/nice/var/logs" || ref_config["path.run"] != "/home/nice/var/run" || ref_config["host.hostname"] != "jelly01.niceprivate.com" {
		t.Logf("wrong ref_config: %#v", ref_config)
		t.Fail()
	}
}

func Test_replae_with_ref(t *testing.T) {
	if err := init_ref_config("ut-conf"); err != nil {
		t.Fatalf("parse ref config failed: %s", err)
	}

	n_val := replace_with_ref("HELLo {$path.var} {$host.hostname} {$haha.xxx} {$haha} llll")
	if n_val != "HELLo /home/nice/var jelly01.niceprivate.com {$haha.xxx} {$haha} llll" {
		t.Logf("replace_with_ref failed: %s", n_val)
		t.Fail()
	}
}
