package os

import (
	cook_util "gitlab.niceprivate.com/golang/cook/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

func Write_pid(pid_fname string) error {
	var err error
	// write pid file
	if len(pid_fname) > 0 {
		if err = ioutil.WriteFile(pid_fname, []byte(strconv.FormatInt(int64(syscall.Getpid()), 10)), 0755); err != nil {
			if !cook_util.Err_NoSuchFileOrDir(err) {
				return err
			}
			if err = os.MkdirAll(filepath.Dir(pid_fname), 0755); err != nil {
				return err
			}
			if err = ioutil.WriteFile(pid_fname, []byte(strconv.FormatInt(int64(syscall.Getpid()), 10)), 0755); err != nil {
				return err
			}
		}
	}

	return nil
}

func Read_pid(pid_fname string) int {
	var (
		pid_fp  *os.File
		content []byte
		pid     int64
		err     error
	)
	if pid_fp, err = os.Open(pid_fname); err != nil {
		return 0
	}
	if content, err = ioutil.ReadAll(pid_fp); err != nil {
		return 0
	}
	if pid, err = strconv.ParseInt(string(content), 10, 64); err != nil {
		return 0
	}

	return int(pid)
}
