package io

import (
	"fmt"
	"net"
	"testing"
	"time"
)

var (
	pR, pW net.Conn
)

func init() {
	var err error

	l, err := net.Listen("tcp", "127.0.0.1:58888")
	if err != nil {
		panic(fmt.Sprintf("listen error: %s", err))
	}
	go func() {
		pR, err = l.Accept()
	}()
	pW, err = net.Dial("tcp", "127.0.0.1:58888")
	if err != nil {
		panic(fmt.Sprintf("dial error: %s", err))
	}
	time.Sleep(1e8)
}

func makeBuffer_net(n int, v byte) []byte {
	buffer := make([]byte, n)
	for i := range buffer {
		buffer[i] = v
	}
	return buffer
}

func TestReadNNormal_net(t *testing.T) {
	done := make(chan string)

	go func(done chan string) {
		begin := time.Now()
		msg, err := ReadN_Net(pR, 100, 150*time.Millisecond)
		if err != nil {
			t.Logf("unexpected error: %s", err)
			t.Fail()
		}
		if time.Now().Sub(begin) < 75*time.Millisecond || time.Now().Sub(begin) > 125*time.Millisecond {
			t.Logf("read used time wrong: want = about 100ms, used = %s", time.Now().Sub(begin))
			t.Fail()
		}
		for i, v := range msg {
			if v != '1' {
				t.Logf("message error at: %d. value: %d", i, v)
				t.Fail()
			}
		}
		done <- "reader"
	}(done)
	go func(done chan string) {
		for i := 0; i < 100; i++ {
			pW.Write([]byte("1"))
			time.Sleep(1e6)
		}
		done <- "writer"
	}(done)
	<-done
	<-done
}

func TestReadNTimeout_net(t *testing.T) {
	done := make(chan string)

	go func(done chan string) {
		begin := time.Now()
		_, err := ReadN_Net(pR, 100, 100*time.Millisecond)
		if time.Now().Sub(begin) < 75*time.Millisecond || time.Now().Sub(begin) > 125*time.Millisecond {
			t.Logf("read used time wrong: want = about 1s, used = %s", time.Now().Sub(begin))
			t.Fail()
		}
		if err == nil || !err.(net.Error).Timeout() {
			t.Logf("unexpected error: %s", err)
			t.Fail()
		}
		done <- "reader"
	}(done)
	go func(done chan string) {
		for i := 0; i < 40; i++ {
			pW.SetWriteDeadline(time.Now().Add(time.Millisecond * 100))
			_, err := pW.Write([]byte("1"))
			if err != nil {
				t.Logf("%s", err)
				t.Fail()
				break
			}
			time.Sleep(1e6)
		}
		time.Sleep(1e8)
		done <- "writer"
	}(done)
	<-done
	<-done
}
