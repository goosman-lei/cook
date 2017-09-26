package daemon

/*
maintence pipe for communicate with parent or child process (for daemon or upgrae)
*/

import (
	cook_log "cook/log"
	"fmt"
	"os"
	"os/exec"
)

var (
	parent_pipe_r *os.File
	parent_pipe_w *os.File
	child_pipe_r  *os.File
	child_pipe_w  *os.File

	upgrade_child_cmd *exec.Cmd
)

const (
	PIPE_C_TO_P_COMPLETE_START = uint8(1)

	PIPE_C_TO_P_SUCC = uint8(2)
	PIPE_C_TO_P_FAIL = uint8(3)
)

func init_pipe_with_child() (err error) {
	child_pipe_r, child_pipe_w, err = os.Pipe()
	return
}

func init_pipe_with_parent() {
	parent_pipe_r = os.NewFile(uintptr(3), "parent_pipe_r")
	parent_pipe_w = os.NewFile(uintptr(4), "parent_pipe_w")
}

func prepend_child_pipe_to_extra_files(extra_files []*os.File) []*os.File {
	child_pipes := []*os.File{child_pipe_r, child_pipe_w}
	if extra_files == nil {
		return child_pipes
	} else {
		return append(child_pipes, extra_files...)
	}
}

func wait_for_child_complete_start() (err error) {
	var buf []byte = make([]byte, 1)

	cook_log.Infof("wait-daemon-start begin")
	if _, err = child_pipe_r.Read(buf); err != nil {
		cook_log.Warnf("wait-daemon-start end with error: %s", err)
		return
	}

	if uint8(buf[0]) != PIPE_C_TO_P_COMPLETE_START {
		err = fmt.Errorf("unkonw pipe message %d", uint8(buf[0]))
		cook_log.Warnf("wait-daemon-start end with error: %s", err)
	} else {
		cook_log.Infof("wait-daemon-start end")
	}
	return
}

func notify_to_parent_complete_start() (err error) {
	cook_log.Infof("notify-parent-started begin")
	var buf []byte = []byte{PIPE_C_TO_P_COMPLETE_START}
	_, err = parent_pipe_w.Write(buf)

	if err == nil {
		cook_log.Infof("notify-parent-started end")
	} else {
		cook_log.Warnf("notify-parent-started end with error: %s", err)
	}
	return
}

func Wait_upgrade_done() bool {
	defer func() {
		if upgrade_child_cmd != nil {
			if err := upgrade_child_cmd.Wait(); err != nil {
				cook_log.Warnf("wait upgrade child failed: %s", err)
			}
			upgrade_child_cmd = nil
		}
	}()
	cook_log.Infof("wait-upgrade-done begin")
	var (
		buf []byte = make([]byte, 1)
		err error
	)
	if _, err = child_pipe_r.Read(buf); err != nil {
		cook_log.Infof("wait-upgrade-done end with error: %s", err)
		return false
	}

	if uint8(buf[0]) == PIPE_C_TO_P_SUCC {
		cook_log.Infof("wait-upgrade-done end with success")
		return true
	} else {
		cook_log.Warnf("wait-upgrade-done end with failure")
		return false
	}
}

func Notify_upgrade_done(success bool) (err error) {
	cook_log.Infof("notify-upgrade-done begin")
	var buf []byte = make([]byte, 1)
	if success {
		buf[0] = PIPE_C_TO_P_SUCC
	} else {
		buf[0] = PIPE_C_TO_P_FAIL
	}
	_, err = parent_pipe_w.Write(buf)

	if err != nil {
		cook_log.Warnf("notify-upgrade-done end with error: %s", err)
		return
	}
	if success {
		cook_log.Infof("notify-upgrade-done end with success")
	} else {
		cook_log.Warnf("notify-upgrade-done end with failure")
	}
	return
}
