package mserver

import (
	"crypto/md5"
	"fmt"
	cook_log "gitlab.niceprivate.com/golang/cook/log"
	cook_util "gitlab.niceprivate.com/golang/cook/util"
	"net"
	"sync"
	"time"
)

type MClient struct {
	Addr       net.Addr
	Sn_literal string
	Sn_md5     string
	Log_prefix string

	OnlineTime     time.Time
	LastActiveTime time.Time

	Svr  *MServer
	Conn *net.TCPConn

	msgQueue chan interface{}

	wg *sync.WaitGroup

	ExtraData *cook_util.CMap
}

func NewMClient(conn *net.TCPConn, svr *MServer) *MClient {
	sn_literal := conn.RemoteAddr().String() + "->" + conn.LocalAddr().String()
	return &MClient{
		Addr:       conn.RemoteAddr(),
		Sn_literal: sn_literal,
		Sn_md5:     fmt.Sprintf("%X", md5.Sum([]byte(sn_literal))),
		Log_prefix: "[" + sn_literal + "]",

		OnlineTime:     time.Now(),
		LastActiveTime: time.Now(),

		Svr:  svr,
		Conn: conn,

		msgQueue: make(chan interface{}, svr.Cfg.ClientMsgQueueSize),

		wg: new(sync.WaitGroup),

		ExtraData: cook_util.NewCMap(),
	}
}

func (c *MClient) run() {
	c.Svr.wg.Add(1)
	c.Debugf("server wait-group +1: MClient.run()#begin")
	defer func() {
		if r := recover(); r != nil {
			// if panic in c.Svr.Ops.ClientOnline, we must cancel waitgroup manual
			c.Svr.wg.Done()
			c.Debugf("server wait-group -1: MClient.run()#recover")
			c.Warnf("panic-client-run: %q", r)
		}
	}()
	// clear bad client on this port
	if oc, e := c.Svr.ClientMap.Get(c.Sn_literal); e {
		oc.(*MClient).Offline()
	}

	if err := c.Svr.Ops.ClientOnline(c); err != nil {
		return
	}

	// must set clientmap before sendloop and recvloop. if else, maybe offline ahead of clientmap.set, it will consider already offline
	c.Svr.ClientMap.Set(c.Sn_literal, c)

	c.wg.Add(2) // send loop and recv loop
	c.Debugf("client wait-group +2: MClient.run()#send_and_recv_loop")
	go c.sendLoop()

	c.recvLoop()
}

func (c *MClient) sendLoop() {
	defer func() {
		c.wg.Done()
		c.Debugf("client wait-group -1: MClient.sendLoop()#defer")
		c.Offline()
		if r := recover(); r != nil {
			c.Warnf("panic-sendloop: %q", r)
		}
	}()
	var (
		msg interface{}
		ok  bool
		err error
	)
SendHandlerLoop:
	for {
		select {
		case msg, ok = <-c.msgQueue:
			// channel closed
			if !ok {
				break SendHandlerLoop
			}
			if err = c.Svr.Ops.SendMsg(c, msg); err != nil {
				c.Offline()
				break SendHandlerLoop
			}
		}
	}
}

func (c *MClient) recvLoop() {
	defer func() {
		c.wg.Done()
		c.Debugf("client wait-group -1: MClient.recvLoop()#defer")
		c.Offline()
		if r := recover(); r != nil {
			c.Warnf("panic-recvloop: %q", r)
		}
	}()

	for {
		if err := c.Svr.Ops.RecvMsg(c); err != nil {
			break
		}
	}
}

func (c *MClient) Offline() {
	defer func() {
		if r := recover(); r != nil {
			c.Warnf("panic-offline: %q", r)
			c.Debugf("server wait-group -1: MClient.Offline()#defer")
			c.Svr.wg.Done()
		}
	}()
	if !c.Svr.ClientMap.CheckAndErase(c.Sn_literal, func(oc interface{}) bool {
		return oc.(*MClient) == c
	}) {
		// already offline
		return
	}

	c.Svr.Ops.ClientOffline(c)
	close(c.msgQueue)
	c.Conn.Close()
	c.Debugf("wait client done: begin")
	c.wg.Wait()
	c.Debugf("wait client done: end")
	c.Debugf("server wait-group -1: MClient.Offline()#end")
	c.Svr.wg.Done()
}

func (c *MClient) Send(msg interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	c.msgQueue <- msg
	return nil
}

func (c *MClient) Debugf(f string, argv ...interface{}) {
	cook_log.Debugf(c.Log_prefix+f, argv...)
}

func (c *MClient) Infof(f string, argv ...interface{}) {
	cook_log.Infof(c.Log_prefix+f, argv...)
}

func (c *MClient) Warnf(f string, argv ...interface{}) {
	cook_log.Warnf(c.Log_prefix+f, argv...)
}

func (c *MClient) Fatalf(f string, argv ...interface{}) {
	cook_log.Fatalf(c.Log_prefix+f, argv...)
}
