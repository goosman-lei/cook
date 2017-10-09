package log

import (
	"fmt"
	cook_util "gitlab.niceprivate.com/golang/cook/util"
	"log"
	"os"
	"syscall"
)

var (
	Ldebug *log.Logger = log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds)
	Linfo  *log.Logger = log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds)
	Lwarn  *log.Logger = log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds)
	Lfatal *log.Logger = log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds)

	pid  int = syscall.Getpid()
	ppid int = syscall.Getppid()

	log_prefix_debug = fmt.Sprintf("debug %d ", pid)
	log_prefix_info  = fmt.Sprintf("info %d ", pid)
	log_prefix_warn  = fmt.Sprintf("warn %d ", pid)
	log_prefix_fatal = fmt.Sprintf("fatal %d ", pid)
)

func SetLogPath(path string) (err error) {
	var logFp *os.File
	if logFp, err = os.OpenFile(path+"/overmind-gate.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755); err == nil {
		goto Set
	}

	if !cook_util.Err_NoSuchFileOrDir(err) {
		return
	}

	if err = os.MkdirAll(path, 0755); err != nil {
		return
	}

	if logFp, err = os.OpenFile(path+"/overmind-gate.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755); err != nil {
		return
	}

Set:
	Ldebug = log.New(logFp, "", log.Ldate|log.Lmicroseconds)
	Linfo = log.New(logFp, "", log.Ldate|log.Lmicroseconds)
	Lwarn = log.New(logFp, "", log.Ldate|log.Lmicroseconds)
	Lfatal = log.New(logFp, "", log.Ldate|log.Lmicroseconds)
	return
}

func Debugf(f string, argv ...interface{}) {
	Ldebug.Printf(log_prefix_debug+f+"\n", argv...)
}

func Infof(f string, argv ...interface{}) {
	Linfo.Printf(log_prefix_info+f+"\n", argv...)
}

func Warnf(f string, argv ...interface{}) {
	Lwarn.Printf(log_prefix_warn+f+"\n", argv...)
}

func Fatalf(f string, argv ...interface{}) {
	Lfatal.Printf(log_prefix_fatal+f+"\n", argv...)
}
