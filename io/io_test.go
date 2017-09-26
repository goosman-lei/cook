package io

import (
	"io"
	"testing"
	"time"
)

func makeBuffer(n int, v byte) []byte {
	buffer := make([]byte, n)
	for i := range buffer {
		buffer[i] = v
	}
	return buffer
}

func TestReadNNormal(t *testing.T) {
	pR, pW := io.Pipe()
	done := make(chan string)

	go func(done chan string) {
		msg, err := ReadN(pR, 100)
		if err != nil {
			t.Logf("unexpected error: %s", err)
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
		}
		done <- "writer"
	}(done)
	<-done
	<-done
}

func TestReadNEOF(t *testing.T) {
	pR, pW := io.Pipe()
	done := make(chan string)

	go func(done chan string) {
		msg, err := ReadN(pR, 100)
		if err != io.EOF {
			t.Logf("unexpected error: %s. already readed message: %s", err, msg)
			t.Fail()
		}
		done <- "reader"
	}(done)
	go func(done chan string) {
		for i := 0; i < 10; i++ {
			pW.Write([]byte("1"))
		}
		pW.Close()
		done <- "writer"
	}(done)
	<-done
	<-done
}

func TestReadNOnSlowWriter(t *testing.T) {
	pR, pW := io.Pipe()
	done := make(chan string)

	go func(done chan string) {
		begin := time.Now()
		msg, err := ReadN(pR, 100)
		end := time.Now()
		if err != nil {
			t.Logf("unexpected error: %s", err)
			t.Fail()
		}
		for i, v := range msg {
			if v != '1' {
				t.Logf("message error at: %d. value: %d", i, v)
				t.Fail()
			}
		}
		if end.Sub(begin) < 1e7 {
			t.Logf("used time less than 100ms: %s", end.Sub(begin))
			t.Fail()
		}
		done <- "reader"
	}(done)
	go func(done chan string) {
		for i := 0; i < 100; i++ {
			pW.Write([]byte("1"))
			time.Sleep(1e5)
		}
		done <- "writer"
	}(done)
	<-done
	<-done
}

func TestWriteNNormal(t *testing.T) {
	pR, pW := io.Pipe()
	done := make(chan string)

	go func(done chan string) {
		buffer := makeBuffer(100, 0)
		nRead, err := pR.Read(buffer[:])
		for i, v := range buffer {
			if v != '1' {
				t.Logf("message error at: %d. value: %d", i, v)
				t.Fail()
			}
		}
		if err != nil {
			t.Logf("unexpected read error: %s", err)
			t.Fail()
		}
		if nRead != 100 {
			t.Logf("message length unexpected: %d", nRead)
			t.Fail()
		}
		done <- "reader"
	}(done)
	go func(done chan string) {
		buffer := makeBuffer(100, '1')
		err := WriteN(pW, buffer[:])
		if err != nil {
			t.Logf("unexpected write error: %s", err)
			t.Fail()
		}
		done <- "writer"
	}(done)
	<-done
	<-done
}

func TestWriteNEOF(t *testing.T) {
	pR, pW := io.Pipe()
	done := make(chan string)

	go func(done chan string) {
		buffer := makeBuffer(10, 0)
		pR.Read(buffer[:])
		pR.Close()
		done <- "reader"
	}(done)
	go func(done chan string) {
		buffer := makeBuffer(100, '1')
		err := WriteN(pW, buffer[:])
		if err != io.ErrClosedPipe {
			t.Logf("expected ErrClosedPipe, but: %s", err)
			t.Fail()
		}
		done <- "writer"
	}(done)
	<-done
	<-done
}

func TestWriteNSlowReader(t *testing.T) {
	pR, pW := io.Pipe()
	done := make(chan string)

	go func(done chan string) {
		buffer := makeBuffer(1, 0)
		nRead := 0
		for ; nRead < 100; nRead++ {
			nTmp, err := pR.Read(buffer[:])
			if err != nil {
				t.Logf("unexpected read error: %s", err)
				t.Fail()
			}
			if nTmp != 1 {
				t.Logf("unexpected read length: %s", nTmp)
				t.Fail()
			}
			time.Sleep(1e5)
		}
		if nRead != 100 {
			t.Logf("message length unexpected: %d", nRead)
			t.Fail()
		}
		pR.Close()
		done <- "reader"
	}(done)
	go func(done chan string) {
		buffer := makeBuffer(101, '1')
		begin := time.Now()
		err := WriteN(pW, buffer[:])
		end := time.Now()
		if end.Sub(begin) < 1e7 {
			t.Logf("used time less than 100ms: %s", end.Sub(begin))
			t.Fail()
		}
		if err != io.ErrClosedPipe {
			t.Logf("expected ErrClosedPipe, but: %s", err)
			t.Fail()
		}
		done <- "writer"
	}(done)
	<-done
	<-done
}
