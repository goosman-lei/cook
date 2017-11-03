package daemon

/*
references:
	* http://software.clapper.org/daemonize/daemonize.html
	* https://github.com/bmc/daemonize/blob/master/daemon.c
	* https://www.socketloop.com/tutorials/golang-daemonizing-a-simple-web-server-process-example
*/

import (
	cook_log "gitlab.niceprivate.com/golang/cook/log"
	cook_util "gitlab.niceprivate.com/golang/cook/util"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

const (
	ENV_NAME_DAEMON_STAGE   = "COOK_DAEMON_STAGE"
	ENV_NAME_LIVING_UPGRADE = "COOK_LIVING_UPGRADE"

	STAGE_PARENT = "parent"
	STAGE_DAEMON = "daemon"
)

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

func CleanEnviron() []string {
	r_env := []string{}
	for _, e := range os.Environ() {
		switch {
		case strings.HasPrefix(e, ENV_NAME_DAEMON_STAGE+"="):
		case strings.HasPrefix(e, ENV_NAME_LIVING_UPGRADE+"="):
		default:
			r_env = append(r_env, e)
		}
	}
	return r_env
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
		Args:       os.Args,
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

	cook_log.Infof("daemon command: %s %s . env_living_upgrade: %s", cook_util.ExecFile(), strings.Join(os.Args[1:], " "), e_string)
	cmd = exec.Cmd{
		Path:       cook_util.ExecFile(),
		Args:       os.Args,
		Env:        os.Environ(),
		Dir:        "/",
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
