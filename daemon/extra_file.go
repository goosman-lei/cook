package daemon

import (
	"bytes"
	"fmt"
	cook_log "gitlab.niceprivate.com/golang/cook/log"
	"net"
	"os"
	"regexp"
	"strconv"
	"sync"
	"syscall"
)

var (
	inherit_label_pattern     *regexp.Regexp = regexp.MustCompile("^[\\w-]+$")
	inherit_env_match_pattern *regexp.Regexp = regexp.MustCompile("^\\s*([\\w-]+)\\s*:\\s*(\\d+)(?:\\s*;\\s*([\\w-]+)\\s*:\\s*(\\d+))*\\s*$")
	inherit_env_find_pattern  *regexp.Regexp = regexp.MustCompile("([\\w-]+)\\s*:\\s*(\\d+)")

	env_daemon_stage   string = os.Getenv(ENV_NAME_DAEMON_STAGE)
	env_living_upgrade string = os.Getenv(ENV_NAME_LIVING_UPGRADE)

	global_extra_mutex   *sync.Mutex           = new(sync.Mutex)
	global_extra_files   map[string][]*os.File = make(map[string][]*os.File)
	global_inherit_files map[string][]*os.File = make(map[string][]*os.File)
)

func init() {
	var (
		matches [][]string
		match   []string
		label   string
		fd_num  int64
		max_fd  int = 4 // stdin, stdout, stderr, pipe_r, pipe_w
		cur_fd  int
	)

	// if upgrade, must have valid inherit file
	if In_upgrade() {
		matches = inherit_env_find_pattern.FindAllStringSubmatch(env_living_upgrade, -1)
		for _, match = range matches {
			label = match[1]
			// if match pattern, fd_num must be int, so ignore error check
			fd_num, _ = strconv.ParseInt(match[2], 10, 64)
			global_inherit_files[label] = make([]*os.File, 0, fd_num)
			for cur_fd = max_fd + 1; fd_num > 0; cur_fd++ {
				global_inherit_files[label] = append(global_inherit_files[label], os.NewFile(uintptr(cur_fd), label))
				cook_log.Infof("in upgrade inherit file: cur_fd = %d label = %s", cur_fd, label)
				fd_num--
				max_fd = cur_fd
			}
		}
	}
}

func In_daemon() bool {
	return env_daemon_stage == STAGE_DAEMON
}

func In_upgrade() bool {
	return inherit_env_match_pattern.MatchString(env_living_upgrade)
}

func AddInherit_TCPListener(label string, l *net.TCPListener, non_blocking bool) error {
	var (
		file   *os.File
		err    error
		exists bool
	)
	if !inherit_label_pattern.MatchString(label) {
		return fmt.Errorf("labe must only can be [\\w-]")
	}
	if file, err = l.File(); err != nil {
		return err
	}

	// (*net.TCPListener).File() will set the fd into blocking mode, it will cause accept timeout ineffective
	// so we can restore it to non-blocking mode
	if non_blocking {
		if err = syscall.SetNonblock(int(file.Fd()), true); err != nil {
			return err
		}
	}

	global_extra_mutex.Lock()
	defer global_extra_mutex.Unlock()

	if _, exists = global_extra_files[label]; !exists {
		global_extra_files[label] = make([]*os.File, 0, 1)
	}

	global_extra_files[label] = append(global_extra_files[label], file)
	return nil
}

func AddInherit_File(label string, file *os.File) error {
	if !inherit_label_pattern.MatchString(label) {
		return fmt.Errorf("labe must only can be [\\w-]")
	}

	global_extra_mutex.Lock()
	defer global_extra_mutex.Unlock()

	if _, exists := global_extra_files[label]; !exists {
		global_extra_files[label] = make([]*os.File, 0, 1)
	}

	global_extra_files[label] = append(global_extra_files[label], file)
	return nil
}

func Inherit_Files(label string) []*os.File {
	if files, exists := global_inherit_files[label]; exists {
		return files
	} else {
		return []*os.File{}
	}
}
func Inherit_TCPListeners(label string) ([]*net.TCPListener, error) {
	defer func() {
		if r := recover(); r != nil {
			cook_log.Infof("inherit panic: %q", r)
		}
	}()
	var (
		ok          bool
		err         error
		file        *os.File
		files       []*os.File
		exists      bool
		listener    net.Listener
		tcpListener *net.TCPListener
		listeners   []*net.TCPListener = make([]*net.TCPListener, 0)
	)
	if files, exists = global_inherit_files[label]; exists {
		for _, file = range files {
			if listener, err = net.FileListener(file); err != nil {
				return nil, err
			}
			if tcpListener, ok = listener.(*net.TCPListener); !ok {
				return nil, fmt.Errorf("not tcp listener")
			}
			listeners = append(listeners, tcpListener)
		}
		return listeners, nil
	} else {
		return []*net.TCPListener{}, nil
	}
}

func Restore_inherit_files() {
	global_extra_mutex.Lock()
	defer global_extra_mutex.Unlock()
	global_extra_files = make(map[string][]*os.File)
}

func extract_living_upgrade_info() ([]*os.File, string) {
	var (
		is_first     bool          = true
		export_files []*os.File    = make([]*os.File, 0)
		buf          *bytes.Buffer = new(bytes.Buffer)
	)

	global_extra_mutex.Lock()
	defer global_extra_mutex.Unlock()
	for label, files := range global_extra_files {
		if is_first {
			is_first = false
			fmt.Fprintf(buf, "%s:%d", label, len(files))
		} else {
			fmt.Fprintf(buf, ";%s:%d", label, len(files))
		}

		export_files = append(export_files, files...)
	}

	return export_files, buf.String()
}

func auto_inherit() {
	for label, files := range global_inherit_files {
		for _, file := range files {
			AddInherit_File(label, file)
		}
	}
}
