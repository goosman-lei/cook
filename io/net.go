package io

import (
	"fmt"
	"net"
	"time"
)

type TimeoutError struct {
	timeout time.Duration
	usetime time.Duration
}

func (e TimeoutError) Error() string {
	return fmt.Sprintf("timeout: usetime = %s timeout: %s", e.usetime, e.timeout)
}
func (e TimeoutError) Timeout() bool   { return true }
func (e TimeoutError) Temporary() bool { return false }

func NewTimeoutError(conn net.Conn, op string, usetime, timeout time.Duration) error {
	return &net.OpError{
		Op:     op,
		Net:    conn.LocalAddr().Network(),
		Source: conn.LocalAddr(),
		Addr:   conn.RemoteAddr(),
		Err: TimeoutError{
			timeout: timeout,
			usetime: usetime,
		},
	}
}

// optimize for gc on large number client recv timeout at same time, it will stop the world if use many local variable in function
type local_scope struct {
	begin     time.Time
	deadline  time.Time
	nRW       int
	tRW       int
	ok        bool
	nErr      net.Error
	nExpected int
}

func ReadN_Net(conn net.Conn, nExpected int, timeout time.Duration) (buffer []byte, err error) {
	ls := local_scope{
		begin:    time.Now(),
		deadline: time.Now().Add(timeout),
		nRW:      0,
		tRW:      0,
	}
	buffer = make([]byte, nExpected)

	for ls.nRW < nExpected {
		conn.SetReadDeadline(ls.deadline)
		ls.tRW, err = conn.Read(buffer[ls.nRW:])
		// have no error, add Read num
		if err == nil {
			ls.nRW += ls.tRW
			continue
		}
		// timeout check
		if ls.deadline.Sub(time.Now()).Nanoseconds() <= 0 {
			err = NewTimeoutError(conn, "read", ls.deadline.Sub(ls.begin), timeout)
			return
		}
		// no net error, or no temporary error
		if ls.nErr, ls.ok = err.(net.Error); !ls.ok || !ls.nErr.Temporary() {
			return
		}
		// temporary net error, will continue reading
	}
	return

}

func ReadToBuffer_Net(conn net.Conn, buffer []byte, timeout time.Duration) (err error) {
	ls := local_scope{
		begin:     time.Now(),
		deadline:  time.Now().Add(timeout),
		nExpected: len(buffer),
		nRW:       0,
		tRW:       0,
	}

	for ls.nRW < ls.nExpected {
		conn.SetReadDeadline(ls.deadline)
		ls.tRW, err = conn.Read(buffer[ls.nRW:])
		// have no error, add Read num
		if err == nil {
			ls.nRW += ls.tRW
			continue
		}
		// timeout check
		if ls.deadline.Sub(time.Now()).Nanoseconds() <= 0 {
			err = NewTimeoutError(conn, "read", ls.deadline.Sub(ls.begin), timeout)
			return
		}
		// no net error, or no temporary error
		if ls.nErr, ls.ok = err.(net.Error); !ls.ok || !ls.nErr.Temporary() {
			return
		}
		// temporary net error, will continue reading
	}
	return nil
}

func WriteN_Net(conn net.Conn, buffer []byte, timeout time.Duration) (err error) {
	ls := local_scope{
		begin:     time.Now(),
		deadline:  time.Now().Add(timeout),
		nExpected: len(buffer),
		nRW:       0,
		tRW:       0,
	}

	for ls.nRW < ls.nExpected {
		conn.SetWriteDeadline(ls.deadline)
		ls.tRW, err = conn.Write(buffer[ls.nRW:])
		// have no error, add Read num
		if err == nil {
			ls.nRW += ls.tRW
			continue
		}
		// timeout check
		if ls.deadline.Sub(time.Now()).Nanoseconds() <= 0 {
			return NewTimeoutError(conn, "write", ls.deadline.Sub(ls.begin), timeout)
		}
		// no net error, or no temporary error
		if ls.nErr, ls.ok = err.(net.Error); !ls.ok || !ls.nErr.Temporary() {
			return err
		}
		// temporary net error, will continue reading
	}
	return nil

}
