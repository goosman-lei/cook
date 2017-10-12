package log

import (
	"io/ioutil"
	"testing"
	"time"
)

func init() {
	Ldebug.SetOutput(ioutil.Discard)
	Linfo.SetOutput(ioutil.Discard)
	Lwarn.SetOutput(ioutil.Discard)
	Lfatal.SetOutput(ioutil.Discard)
}

func TestCase_log_file_rotate(t *testing.T) {
	t.Skip()
	Debugf("Hi: %s", "Jack")
	SetupLog([]LogConf{
		LogConf{Level: LEVEL_DEBUG, Fpath: "../debug.log-*-*-*-*-*"},
		LogConf{Level: LEVEL_INFO, Fpath: "../info.log-*-*-*-*-*"},
		LogConf{Level: LEVEL_WARN, Fpath: "../warn.log-*-*-*-*-*"},
		LogConf{Level: LEVEL_FATAL, Fpath: "../fatal.log-*-*-*-*-*"},
	})

	for i := 0; i < 10; i++ {
		Debugf("log index: %d", i)
		Infof("log index: %d", i)
		Warnf("log index: %d", i)
		Fatalf("log index: %d", i)
		time.Sleep(time.Minute)
	}
	//t.Fail()
}

func Benchmark_Debugf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Debugf("Hi: %s", "Jack")
	}
}
