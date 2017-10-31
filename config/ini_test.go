package config

import (
	"testing"
)

func Test_Init_open_dir(t *testing.T) {
	iniFp, err := Ini_open_dir("ut-conf")
	if err != nil {
		t.Logf("open dir failed: %s", err)
		t.Fail()
	}

	sec, err := iniFp.GetSection("Log")
	if err != nil {
		t.Logf("get section Log failed: %s", err)
		t.Fail()
	}

	if Ini_direct_get_key(sec, "Log", "DebugPath").MustString("") != "/home/nice/var/logs/debug.log" {
		t.Logf("wrong Log.DebugPath: %s", Ini_direct_get_key(sec, "Log", "DebugPath").MustString(""))
		t.Fail()
	}

	if Ini_direct_get_key(sec, "Log", "DataPath").MustString("") != "/home/nice/var/data/data.log" {
		t.Logf("wrong Log.DebugPath: %s", Ini_direct_get_key(sec, "Log", "DataPath").MustString(""))
		t.Fail()
	}
}
