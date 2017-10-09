package mserver

import (
	cook_daemon "gitlab.niceprivate.com/golang/cook/daemon"
	cook_log "gitlab.niceprivate.com/golang/cook/log"
	cook_util "gitlab.niceprivate.com/golang/cook/util"
	"net"
	"time"
)

type MListeners struct {
	Listeners []*MListener
	Svr       *MServer
}
type MListener struct {
	literalAddr string
	NetAddr     net.Addr
	Listener    *net.TCPListener
	Svr         *MServer
	Ls          *MListeners

	Log_prefix string
}

func (s *MServer) initListeners() {
	s.Ls = &MListeners{}
	s.Ls.Listeners = make([]*MListener, len(s.Cfg.ListenAddrs))
	s.Ls.Svr = s
	for i, addr := range s.Cfg.ListenAddrs {
		s.Ls.Listeners[i] = &MListener{
			Svr:         s,
			Ls:          s.Ls,
			literalAddr: addr,
			Log_prefix:  "[" + s.SN + "#" + addr + "]",
		}
	}
}

func (ls *MListeners) listen() error {
	if cook_daemon.In_upgrade() {
		return ls.upgradeListen()
	} else {
		return ls.createListen()
	}
}

func (ls *MListeners) createListen() error {
	var (
		err      error
		listener *net.TCPListener
		tcpAddr  *net.TCPAddr
	)
	for _, l := range ls.Listeners {
		if tcpAddr, err = net.ResolveTCPAddr("tcp", l.literalAddr); err != nil {
			return err
		}

		if listener, err = net.ListenTCP("tcp", tcpAddr); err != nil {
			return err
		}

		ls.Svr.Infof("new listen on %s", listener.Addr().String())

		l.Listener = listener
		l.NetAddr = listener.Addr()
	}

	return nil
}

func (ls *MListeners) upgradeListen() error {
	var (
		tcpAddr           *net.TCPAddr
		exists            bool
		err               error
		listener          *net.TCPListener
		listeners         []*net.TCPListener
		listeners_mapping map[int]*net.TCPListener
	)

	if listeners, err = cook_daemon.Inherit_TCPListeners(ls.Svr.Cfg.UpgradeFileKey); err != nil {
		// inherit error, must notify to parent(which process startup upgrade) this error
		ls.Svr.Fatalf("fetch inherit listeners failed: %s", err)
		return err
	}

	// build mapping
	listeners_mapping = make(map[int]*net.TCPListener)
	for _, listener = range listeners {
		listeners_mapping[listener.Addr().(*net.TCPAddr).Port] = listener
	}

	// restore from inherit or create new listen
	for _, l := range ls.Listeners {
		if tcpAddr, err = net.ResolveTCPAddr("tcp", l.literalAddr); err != nil {
			ls.Svr.Fatalf("resolve tcpaddr(%s) failed: %s", l.literalAddr, err)
			return err
		}

		if listener, exists = listeners_mapping[tcpAddr.Port]; !exists {
			// new listen addr
			if listener, err = net.ListenTCP("tcp", tcpAddr); err != nil {
				ls.Svr.Fatalf("create new listener on %d failed: %s", tcpAddr.Port, err)
				return err
			}
			ls.Svr.Infof("new listen on %s", tcpAddr.Port)
		} else {
			delete(listeners_mapping, tcpAddr.Port)
			ls.Svr.Infof("inherit listen on %s", listener.Addr().String())
		}

		l.Listener = listener
		l.NetAddr = listener.Addr()
	}

	// other inherit file, direct close it
	for addr, l := range listeners_mapping {
		ls.Svr.Infof("inherit from parent, but have not config, close it: %s", addr)
		l.Close()
	}

	return nil
}

func (ls *MListeners) accept() {
	for _, l := range ls.Listeners {
		go l.accept()
	}
}

func (ls *MListeners) shutdown() {
	for _, l := range ls.Listeners {
		l.Listener.Close()
		ls.Svr.Infof("shutdown, close listen: %s", l.NetAddr.String())
	}
}

func (ls *MListeners) upgrade_prepare() error {
	for _, l := range ls.Listeners {
		if err := cook_daemon.AddInherit_TCPListener(ls.Svr.Cfg.UpgradeFileKey, l.Listener, true); err != nil {
			return err
		}
		ls.Svr.Infof("add %s into inherit file list", l.NetAddr.String())
	}
	return nil
}

func (ls *MListeners) upgrade_done() {
	for _, l := range ls.Listeners {
		l.Listener.Close()
		ls.Svr.Infof("upgrade-done, close listen: %s", l.NetAddr.String())
	}
}

// running util server status is not running
func (l *MListener) accept() {
	l.Svr.wg.Add(1)
	defer l.Svr.wg.Done()
	var (
		conn *net.TCPConn
		err  error
	)
	l.Infof("accept() start")
AcceptLoop:
	for {
		select {
		case <-l.Svr.doneCh:
			break AcceptLoop
		default:
		}
		l.Listener.SetDeadline(time.Now().Add(time.Millisecond * 10))
		if conn, err = l.Listener.AcceptTCP(); err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			if cook_util.Err_IsClosed(err) || cook_util.Err_IsBroken(err) {
				break
			}
			l.Warnf("accept failed: %s", err)
			continue
		}
		go NewMClient(conn, l.Svr).run()
	}
	l.Infof("accept() end")
}

func (l *MListener) Debugf(f string, argv ...interface{}) {
	cook_log.Debugf(l.Log_prefix+f, argv...)
}

func (l *MListener) Infof(f string, argv ...interface{}) {
	cook_log.Infof(l.Log_prefix+f, argv...)
}

func (l *MListener) Warnf(f string, argv ...interface{}) {
	cook_log.Warnf(l.Log_prefix+f, argv...)
}

func (l *MListener) Fatalf(f string, argv ...interface{}) {
	cook_log.Fatalf(l.Log_prefix+f, argv...)
}
