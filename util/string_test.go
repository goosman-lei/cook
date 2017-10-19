package util

import (
	"testing"
)

func Test_Hump_to_underline(t *testing.T) {
	if "get_user_settings" != Hump_to_underline("GetUserSettings") {
		t.Logf("faield")
		t.Fail()
	}
	if "get_user_settings" != Hump_to_underline("getUserSettings") {
		t.Logf("faield")
		t.Fail()
	}
	if "get_user_settings" != Hump_to_underline("getUserSettings") {
		t.Logf("faield")
		t.Fail()
	}
	if "get_user_settings" != Hump_to_underline("get_userSettings") {
		t.Logf("faield")
		t.Fail()
	}
}
