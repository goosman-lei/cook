package pool

import (
	cook_util "cook/util"
	"net"
	"sync"
	"testing"
	"time"
)

func startup_server_listener(t *testing.T) {
	var (
		l net.Listener
		e error
		c net.Conn
	)

	if l, e = net.Listen("tcp", ":2021"); e != nil {
		t.Fatalf("listen failed: %s", e)
	}

	go func() {
		for {
			l.(*net.TCPListener).SetDeadline(time.Now().Add(time.Second))
			if c, e = l.Accept(); e != nil {
				t.Fatalf("accept error: %s", e)
				return
			}
			t.Logf("connect from: %s", c.RemoteAddr())
		}
	}()

}

func TestCase_pool_conn(t *testing.T) {
	startup_server_listener(t)

	var (
		pool  *Pool_conn
		wg    *sync.WaitGroup           = new(sync.WaitGroup)
		conns map[string]*Pool_conn_obj = make(map[string]*Pool_conn_obj, 5)
	)

	pool = NewPool_conn(5, func() (net.Conn, error) {
		return net.Dial("tcp", "127.0.0.1:2021")
	})

	for i := 0; i < 5; i++ {
		if conn, err := pool.Get_timeout(time.Millisecond); err != nil {
			t.Fatalf("get connection from pool failed: %s", err)
		} else {
			conns[conn.LocalAddr().String()] = conn
		}
	}
	if conn, err := pool.Get_timeout(time.Millisecond); err != ErrTimeout || conn != nil {
		t.Fatalf("Except timeout error, but have noe: %q %s", conn, err)
	}
	for _, conn := range conns {
		conn.Close()
	}

	for i := 0; i < 100; i++ {
		go func(goid int) {
			wg.Add(1)
			defer wg.Done()
			for j := 0; j < 10; j++ {
				conn, err := pool.Get_timeout(time.Millisecond)
				if (conn == nil && err != ErrTimeout) || (conn != nil && err != nil) {
					t.Fatalf("in goroutine error: %q %s", conn, err)
				}

				if err == nil {
					//t.Logf("goroutine %d success %s", goid, conn.LocalAddr().String())
				} else {
					//t.Logf("goroutine %d timeout", goid)
				}

				if err == nil {
					conn.Close()
				}
			}
		}(i)
	}

	wg.Wait()

	for i := 0; i < 5; i++ {
		if conn, err := pool.Get_timeout(time.Millisecond); err != nil {
			t.Fatalf("get connection from pool failed: %s", err)
		} else {
			if _, exists := conns[conn.LocalAddr().String()]; !exists {
				t.Fatalf("after concurrence running. wrone conn: %s %q", conn.LocalAddr().String(), conns)
			}
			conn.Close()
		}
	}

	if conn, err := pool.Get_timeout(time.Millisecond); err != nil {
		t.Fatalf("get connection from pool failed: %s", err)
	} else {
		pool.Close()
		if err := conn.Conn.SetReadDeadline(time.Now().Add(time.Millisecond)); err != nil {
			t.Fatalf("unexcpet error: %s", err)
		}
		// because pool is closed. so when conn close, it will close underlay net.Conn
		conn.Close()
		if err := conn.Conn.SetReadDeadline(time.Now().Add(time.Millisecond)); err == nil {
			t.Fatalf("excpet use closed socket error, but no error")
		} else if !cook_util.Err_IsClosed(err) {
			t.Fatalf("excpet use closed socket error, but: %s", err)
		}
	}
}
