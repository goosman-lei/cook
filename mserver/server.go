package mserver

import (
	cook_log "gitlab.niceprivate.com/golang/cook/log"
	cook_util "gitlab.niceprivate.com/golang/cook/util"
	"sync"
)

type MSvrConf struct {
	SN          string
	ListenAddrs []string

	ClientMsgQueueSize int

	UpgradeFileKey string
}

type MSvrOps struct {
	ServerStartup  func(*MServer) error
	ServerShutdown func(*MServer) error

	ServerUpgrade func(*MServer) error
	ClientUpgrade func(*MClient) error

	ClientOnline  func(*MClient) error
	ClientOffline func(*MClient) error

	SendMsg func(*MClient, interface{}) error
	RecvMsg func(*MClient) error
}

type MServer struct {
	SN         string
	Cfg        MSvrConf
	Ops        MSvrOps
	Log_prefix string

	Ls *MListeners

	ClientMap *cook_util.CMap

	wg     *sync.WaitGroup
	doneCh chan bool
}

func NewMServer(cfg MSvrConf, ops MSvrOps) *MServer {
	s := &MServer{
		SN:         cfg.SN,
		Cfg:        cfg,
		Ops:        ops,
		Log_prefix: "[" + cfg.SN + "]",

		ClientMap: cook_util.NewCMap(),

		wg:     new(sync.WaitGroup),
		doneCh: make(chan bool),
	}

	s.initListeners()

	return s
}

func (s *MServer) GetListenAddrs() []string {
	return s.Cfg.ListenAddrs
}

func (s *MServer) Startup() (err error) {
	// listening
	if err = s.Ls.listen(); err != nil {
		return
	}

	// application startup
	if err = s.Ops.ServerStartup(s); err != nil {
		s.Fatalf("Ops.ServerStartup() failed: %s", err)
		return
	}

	// accept
	s.Ls.accept()

	return
}

func (s *MServer) Broadcast(msg interface{}) {
	cs := s.ClientMap.Vals()
	for _, c := range cs {
		c.(*MClient).Send(msg)
	}
}

func (s *MServer) Shutdown(reason string) {
	s.Warnf("will shutdown, reason: %s", reason)

	// application shutdown
	if err := s.Ops.ServerShutdown(s); err != nil {
		s.Fatalf("Ops.ServerShutdown() failed: %s", err)
	}

	// stop accept handler
	close(s.doneCh)
	// close listen socket
	s.Ls.shutdown()

	// because s.ClientMap.Iterate() is not safe, if in iteration used cmap' method
	// first, fetch all clients
	cs := s.ClientMap.Vals()
	for _, c := range cs {
		c.(*MClient).Offline()
	}

	s.wg.Wait()
}

func (s *MServer) Upgrade_prepare() error {
	s.Warnf("will upgrade")

	// application upgrade
	if err := s.Ops.ServerUpgrade(s); err != nil {
		s.Fatalf("Ops.ServerUpgrade() failed: %s", err)
		return err
	}

	// prepare inherit listen
	return s.Ls.upgrade_prepare()
}

func (s *MServer) Upgrade_done() {
	// stop accept handler
	close(s.doneCh)
	// close listen socket
	s.Ls.upgrade_done()

	// here, give a chance for application, notify client graceful offline
	// we don't force offline, unless application return error
	cs := s.ClientMap.Vals()
	for _, c := range cs {
		if err := s.Ops.ClientUpgrade(c.(*MClient)); err != nil {
			c.(*MClient).Infof("upgrade failed: %s", err)
			c.(*MClient).Offline()
		}
	}
}

func (s *MServer) Wait() {
	s.wg.Wait()
}

/*
Performance Testing for variable parameter merge
a = fmt.Sprintf(prefix+f, argv...)
    len(argv) == 4:  372 ns/op
    len(argv) == 128: 5509 ns/op
    len(argv) == 1024:  45243 ns/op
a = fmt.Sprintf("["+node+"#"+key+"]"+f, argv...)
    len(argv) == 4:  412 ns/op
    len(argv) == 128: 5886 ns/op
    len(argv) == 1024:  45369 ns/op
a = fmt.Sprintf("[%s#%s]"+f, append([]interface{}{node, key}, argv...))
    len(argv) == 4:  2273 ns/op
    len(argv) == 128: 35116 ns/op
    len(argv) == 1024:  317331 ns/op
*/

func (s *MServer) Debugf(f string, argv ...interface{}) {
	cook_log.Debugf(s.Log_prefix+f, argv...)
}

func (s *MServer) Infof(f string, argv ...interface{}) {
	cook_log.Infof(s.Log_prefix+f, argv...)
}

func (s *MServer) Warnf(f string, argv ...interface{}) {
	cook_log.Warnf(s.Log_prefix+f, argv...)
}

func (s *MServer) Fatalf(f string, argv ...interface{}) {
	cook_log.Fatalf(s.Log_prefix+f, argv...)
}
