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

func ReadN_Net(conn net.Conn, nExpected int, timeout time.Duration) (buffer []byte, err error) {
	var (
		begin    time.Time = time.Now()
		deadline time.Time = begin.Add(timeout)
		nRW      int       = 0
		tRW      int       = 0
		ok       bool
		nErr     net.Error
	)
	buffer = make([]byte, nExpected)

	for nRW < nExpected {
		conn.SetReadDeadline(deadline)
		tRW, err = conn.Read(buffer[nRW:])
		// have no error, add Read num
		if err == nil {
			nRW += tRW
			continue
		}
		// timeout check
		if deadline.Sub(time.Now()).Nanoseconds() <= 0 {
			err = NewTimeoutError(conn, "read", deadline.Sub(begin), timeout)
			return
		}
		// no net error, or no temporary error
		if nErr, ok = err.(net.Error); !ok || !nErr.Temporary() {
			return
		}
		// temporary net error, will continue reading
	}
	return

}

func ReadToBuffer_Net(conn net.Conn, buffer []byte, timeout time.Duration) (err error) {
	var (
		begin     time.Time = time.Now()
		deadline  time.Time = begin.Add(timeout)
		nRW       int       = 0
		nExpected int       = len(buffer)
		tRW       int       = 0
		ok        bool
		nErr      net.Error
	)

	for nRW < nExpected {
		conn.SetReadDeadline(deadline)
		tRW, err = conn.Read(buffer[nRW:])
		// have no error, add Read num
		if err == nil {
			nRW += tRW
			continue
		}
		// timeout check
		if deadline.Sub(time.Now()).Nanoseconds() <= 0 {
			err = NewTimeoutError(conn, "read", deadline.Sub(begin), timeout)
			return
		}
		// no net error, or no temporary error
		if nErr, ok = err.(net.Error); !ok || !nErr.Temporary() {
			return
		}
		// temporary net error, will continue reading
	}
	return nil
}

func WriteN_Net(conn net.Conn, buffer []byte, timeout time.Duration) (err error) {
	var (
		begin     time.Time = time.Now()
		deadline  time.Time = begin.Add(timeout)
		nRW       int       = 0
		nExpected int       = len(buffer)
		tRW       int       = 0
		ok        bool
		nErr      net.Error
	)

	for nRW < nExpected {
		conn.SetWriteDeadline(deadline)
		tRW, err = conn.Write(buffer[nRW:])
		// have no error, add Read num
		if err == nil {
			nRW += tRW
			continue
		}
		// timeout check
		if deadline.Sub(time.Now()).Nanoseconds() <= 0 {
			return NewTimeoutError(conn, "write", deadline.Sub(begin), timeout)
		}
		// no net error, or no temporary error
		if nErr, ok = err.(net.Error); !ok || !nErr.Temporary() {
			return err
		}
		// temporary net error, will continue reading
	}
	return nil
}
