# About daemonize

```go
if is_daemon, err := cook_daemon.Daemonize(); err != nil {
    // error handling
    // here, only in parent
} else if !is_daemon {
    // here is parent process
    // maybe you can do something about clean
    // and immediatly, you should exit it
    os.Exit(0)
}
// ok, here is daemon child process, work at here

// if you make sure, server normal startup, you can write pid file manual
cook_daemon.Write_pid("pid-file-name")

/*
After your program startup, process diagram like this:

6590                parent(immediatly will exit)
|
| ---- 6595         daemon(will running of long time)
*/
```

# About Upgrade

```go
if is_daemon, err := cook_daemon.Daemonize(); err != nil {
    // error handling
    // here, only in parent
} else if !is_daemon {
    // here is parent process
    // maybe you can do something about clean
    // and immediatly, you should exit it
    os.Exit(0)
}
// ok, here is daemon child process, work at here

// here we create listen socket, with living-upgrade check
var (
    listener *net.TCPListener
    listeners []*net.TCPListener
)
if cook_daemon.In_upgrade() {
    // fetch listener from parent process
    listeners, _ = cook_daemon.Inherit_TCPListeners("listen")
    // some other code
} else {
    // create new listener
    listener, _ = net.ListenTCP("tcp", ":2021")
    // some other code
}

// if you make sure, server complete startup, and server running in upgrade mode
// you must notify this to parent
cook_daemon.Notify_upgrade_done()

// if you make sure, server normal startup, you can write pid file manual
cook_daemon.Write_pid("pid-file-name")

// here, we can install a signal handler
// when it receive SIGUSR1, it will execute living-upgrade
signal.Notify(sigCh, syscall.SIGUSR1)
go func() {
    for sig := range sigCh {
        switch sig {
        case syscall.SIGUSR1:
            // add inherit file
            cook_daemon.AddInherit_TCPListener("listen", listener, true)

            if err := cook_daemon.Upgrade(); err != nil {
                // upgrade failed, log it

                // resotre inherit file
                cook_daemon.Restore_inherit_files()

                // other restore works
            } else {
                // some clean work, for example:
                // * close listener
                // * send exit msg to all client which already connected
                // * wait all client exit
                // * close server
                // * and finnaly, exit this process
            }
        }
    }
}()

/*
Notice: in upgrade process, you can run upgrade again.

After your program startup, process diagram like this:

6590                                            parent(immediatly will exit)
|
| ---- 6595                                     daemon(will runing of long time, util itself exit or all client offline after upgrade)
        |
        | ---- 6613                             upgrade-parent(immediatly will exit)
                |
                | ---- 6621                     upgrade-daemon(will runing of long time, util itself exit or all client offline after upgrade)
                        |
                        | ---- 6765             upgrade-parent(immediatly will exit)
                                |
                                | ---- 6773     upgrade-daemon(will runing of long time, util itself exit or all client offline after upgrade)
*/

```

# Demo code for upgrade and daemonize

## code

```go
package main

import (
	cook_daemon "cook/daemon"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	log_fp *os.File = os.Stderr
)

func logf(f string, argv ...interface{}) {
	fmt.Fprintf(log_fp, "%d "+f+"\n", append([]interface{}{syscall.Getpid()}, argv...)...)
}

func main() {
	var (
		err error
	)
	if log_fp, err = os.OpenFile("/tmp/log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755); err != nil {
		logf("open log file failed: %s", err)
		return
	}

	logf("running environ: stage = %s living_upgrade= %s", os.Getenv(cook_daemon.ENV_NAME_DAEMON_STAGE), os.Getenv(cook_daemon.ENV_NAME_LIVING_UPGRADE))

	if in_daemon, err := cook_daemon.Daemonize(); err != nil {
		logf("daemonized failed: %s", err)
		return
	} else if !in_daemon {
		logf("daemonized success. parent will exit")
		os.Exit(0)
	}

	logf("daemonized success. child-daemon begin running")

	if cook_daemon.In_upgrade() {
		upgrade_handler()
	} else {
		normal_handler()
	}
	logf("child-daemon end")
}

func normal_handler() {
	var (
		tcpAddr  *net.TCPAddr
		listener *net.TCPListener
		err      error
	)
	if tcpAddr, err = net.ResolveTCPAddr("tcp", ":2021"); err != nil {
		logf("resolve tcpaddr failed: %s", err)
		return
	}
	if listener, err = net.ListenTCP("tcp", tcpAddr); err != nil {
		logf("create listener failed: %s", err)
	}
	register_signal(listener)

	for {
		listener.SetDeadline(time.Now().Add(time.Second * 10))
		conn, err := listener.Accept()
		if nErr, ok := err.(net.Error); ok && nErr.Timeout() {
			logf("normal listener timeout: %s", time.Now())
			continue
		}
		logf("normal listener accept from: %s", conn.RemoteAddr().String())
	}
}

func upgrade_handler() {
	var (
		listeners []*net.TCPListener
		listener  *net.TCPListener
		err       error
	)
	if listeners, err = cook_daemon.Inherit_TCPListeners("listen"); err != nil {
		cook_daemon.Notify_upgrade_done(false)
		logf("fetch inherit listener failed: %s", err)
		return
	}
	listener = listeners[0]

	cook_daemon.Notify_upgrade_done(true)

	register_signal(listener)

	for {
		listener.SetDeadline(time.Now().Add(time.Second * 10))
		conn, err := listener.Accept()
		if nErr, ok := err.(net.Error); ok && nErr.Timeout() {
			logf("upgrade listener timeout: %s", time.Now())
			continue
		}
		logf("upgrade listener accept from: %s", conn.RemoteAddr().String())
	}
}

func register_signal(listener *net.TCPListener) {
	var (
		sigCh chan os.Signal = make(chan os.Signal)
		sig   os.Signal
		err   error
	)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGPIPE, syscall.SIGCHLD, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGUSR1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logf("panic: %q", r)
			}
		}()
		for sig = range sigCh {
			logf("receive signal: %s", sig)
			switch sig {
			case syscall.SIGUSR1:
				if err = cook_daemon.AddInherit_TCPListener("listen", listener, true); err != nil {
					logf("add listener to remain file list failed: %s", err)
					continue // ignore failed
				}
				logf("add listener to inherit file")
				if err = cook_daemon.Upgrade(); err != nil {
					logf("upgrade error: %s", err)
					cook_daemon.Restore_inherit_files()
					continue // ignore failed
				}
				logf("execute upgrade")
				// here you would do something about clean. or wait client exit
				time.Sleep(1e10)
				os.Exit(0)
			}
		}
	}()
}
```

## command and output

```
$ go build src/main.go
$ ./main

output and interactive command

38585 running environ: stage =  living_upgrade=                     | # parent
38589 running environ: stage = daemon living_upgrade=               | # daemon
38585 daemonized success. parent will exit                          |
38589 daemonized success. child-daemon begin running                |
38589 receive signal: user defined signal 1                         | kill -s SIGUSR1 $(cat /tmp/pid)
38589 add listener to inherit file                                  | # upgrade
38628 running environ: stage = parent living_upgrade= listen:1      | # upgrade-parent
38636 running environ: stage = daemon living_upgrade= listen:1      | # upgrade-daemon
38628 daemonized success. parent will exit                          |
38636 daemonized success. child-daemon begin running                |
38589 execute upgrade                                               |
38636 upgrade listener accept from: 127.0.0.1:32503                 | echo "" | nc 127.0.0.1 2021
38589 normal listener accept from: 127.0.0.1:32507                  | echo "" | nc 127.0.0.1 2021
38636 upgrade listener accept from: 127.0.0.1:32513                 | echo "" | nc 127.0.0.1 2021
38636 upgrade listener accept from: 127.0.0.1:32516                 | echo "" | nc 127.0.0.1 2021
38589 normal listener accept from: 127.0.0.1:32519                  | echo "" | nc 127.0.0.1 2021
38589 normal listener accept from: 127.0.0.1:32524                  | echo "" | nc 127.0.0.1 2021
38589 normal listener accept from: 127.0.0.1:32526                  | echo "" | nc 127.0.0.1 2021
38589 normal listener accept from: 127.0.0.1:32574                  | echo "" | nc 127.0.0.1 2021
38589 normal listener accept from: 127.0.0.1:32593                  | echo "" | nc 127.0.0.1 2021
38589 normal listener accept from: 127.0.0.1:32596                  | echo "" | nc 127.0.0.1 2021
38636 upgrade listener accept from: 127.0.0.1:32598                 | echo "" | nc 127.0.0.1 2021  # after this line, original listen socket complete 10s sleep, process sleep
38636 upgrade listener accept from: 127.0.0.1:32599                 | echo "" | nc 127.0.0.1 2021
38636 upgrade listener accept from: 127.0.0.1:32603                 | echo "" | nc 127.0.0.1 2021
38636 upgrade listener accept from: 127.0.0.1:32604                 | echo "" | nc 127.0.0.1 2021
38636 upgrade listener accept from: 127.0.0.1:32606                 | echo "" | nc 127.0.0.1 2021
38636 upgrade listener accept from: 127.0.0.1:32609                 | echo "" | nc 127.0.0.1 2021
38636 upgrade listener accept from: 127.0.0.1:32610                 | echo "" | nc 127.0.0.1 2021
38636 upgrade listener accept from: 127.0.0.1:32611                 | echo "" | nc 127.0.0.1 2021
38636 upgrade listener accept from: 127.0.0.1:32613                 | echo "" | nc 127.0.0.1 2021
```
