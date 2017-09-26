package daemon

/*
references:
	* http://software.clapper.org/daemonize/daemonize.html
	* https://github.com/bmc/daemonize/blob/master/daemon.c
	* https://www.socketloop.com/tutorials/golang-daemonizing-a-simple-web-server-process-example
*/

import (
	cook_log "cook/log"
	cook_util "cook/util"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

const (
	ENV_NAME_DAEMON_STAGE   = "COOK_DAEMON_STAGE"
	ENV_NAME_LIVING_UPGRADE = "COOK_LIVING_UPGRADE"

	STAGE_PARENT = "parent"
	STAGE_DAEMON = "daemon"
)

func Getpid(pid_fname string) int {
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

// we must run daemonize at first, if you want daemon, because it's not thread safe, but golang use thread manage runtime/gc, an goroutine
func Daemonize() (bool, error) {
	// already is daemon
	if !In_daemon() && syscall.Getppid() == 1 {
		return true, nil
	}

	if !In_daemon() {
		return false, do_parent()
	} else {
		return true, do_daemon()
	}

}

func Upgrade() error {
	defer func() {
		if r := recover(); r != nil {
			cook_log.Infof("panic from upgrade: %q", r)
		}
	}()
	var (
		e_files  []*os.File
		e_string string
		err      error
	)

	// make pipe for child <--> parent
	if err = init_pipe_with_child(); err != nil {
		return err
	}

	e_files, e_string = extract_living_upgrade_info()
	e_files = prepend_child_pipe_to_extra_files(e_files)

	os.Setenv(ENV_NAME_LIVING_UPGRADE, e_string)

	cook_log.Infof("upgrade command: %s %s . env_living_upgrade: %s", cook_util.ExecFile(), strings.Join(os.Args[1:], " "), e_string)
	upgrade_child_cmd = &exec.Cmd{
		Path:       cook_util.ExecFile(),
		Args:       os.Args[1:],
		Env:        os.Environ(),
		ExtraFiles: e_files,
	}

	if err = upgrade_child_cmd.Start(); err != nil {
		return err
	}

	// we must ensure parent is living, before child start success.
	// reference: $GOROOT/src/syscall/exec_linux.go forkAndExecInChild()
	// if parent is already dead, child will kill self
	if err = wait_for_child_complete_start(); err != nil {
		cook_log.Infof("wait upgrade child complete start failed: %s", err)
		return err
	}

	return nil
}
func do_parent() error {
	var (
		e_files  []*os.File
		e_string string
		dev_null *os.File
		err      error
		cmd      exec.Cmd
	)

	// make pipe for child <--> parent
	if err = init_pipe_with_child(); err != nil {
		return err
	}

	if In_upgrade() {
		init_pipe_with_parent()
		notify_to_parent_complete_start()
		auto_inherit()
		e_files, e_string = extract_living_upgrade_info()
		e_files = prepend_child_pipe_to_extra_files(e_files)
	} else {
		e_files, e_string = prepend_child_pipe_to_extra_files(nil), ""
	}

	os.Setenv(ENV_NAME_DAEMON_STAGE, STAGE_DAEMON)
	os.Setenv(ENV_NAME_LIVING_UPGRADE, e_string)

	if dev_null, err = os.OpenFile(os.DevNull, os.O_RDWR, 0755); err != nil {
		return fmt.Errorf("open /dev/null failed: %s", err)
	}

	cook_log.Infof("daemon command: %s %s . env_living_upgrade: %s", cook_util.ExecFile(), strings.Join(os.Args[1:], " "), e_string)
	cmd = exec.Cmd{
		Path:       cook_util.ExecFile(),
		Args:       os.Args[1:],
		Env:        os.Environ(),
		Dir:        "/",
		Stdin:      dev_null,
		Stdout:     dev_null,
		Stderr:     dev_null,
		ExtraFiles: e_files,
		SysProcAttr: &syscall.SysProcAttr{
			Setsid: true,
		},
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	// we must ensure parent is living, before child start success.
	// reference: $GOROOT/src/syscall/exec_linux.go forkAndExecInChild()
	// if parent is already dead, child will kill self
	if err = wait_for_child_complete_start(); err != nil {
		cook_log.Infof("wait daemon complete start failed: %s", err)
	}
	// when upgrade, parent should be closed immediatly, after child running done
	if In_upgrade() {
		Notify_upgrade_done(Wait_upgrade_done())
		os.Exit(0)
	}
	return nil
}

func do_daemon() error {
	var err error
	// here we must have a sync message, notify parent: child startup done
	// reference do_parent()
	init_pipe_with_parent()
	if err = notify_to_parent_complete_start(); err != nil {
		return err
	}

	os.Setenv(ENV_NAME_DAEMON_STAGE, STAGE_PARENT)

	syscall.Umask(0)

	return nil
}

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
