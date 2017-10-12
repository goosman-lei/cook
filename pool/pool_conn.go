package pool

import (
	"net"
	"time"
)

type Pool_conn struct {
	*Pool
}

type Pool_conn_obj struct {
	net.Conn
	pc *Pool_conn
}

func NewPool_conn(c int, factory func() (net.Conn, error)) *Pool_conn {
	pc := &Pool_conn{}
	pc.Pool = new_pool(c, func() (interface{}, error) {
		if c, err := factory(); err != nil {
			return nil, err
		} else {
			return &Pool_conn_obj{
				Conn: c,
				pc:   pc,
			}, nil
		}
	}, func(d interface{}) {
		if c, ok := d.(*Pool_conn_obj); ok {
			c.Conn.Close()
		}
	})

	return pc
}

func (pc *Pool_conn) Get() (*Pool_conn_obj, error) {
	if c, err := pc.get(); err != nil {
		return nil, err
	} else if v, ok := c.(*Pool_conn_obj); !ok {
		return nil, ErrNotWrapper
	} else {
		return v, nil
	}
}

func (pc *Pool_conn) Get_timeout(duration time.Duration) (*Pool_conn_obj, error) {
	if c, err := pc.get_timeout(duration); err != nil {
		return nil, err
	} else if v, ok := c.(*Pool_conn_obj); !ok {
		return nil, ErrNotWrapper
	} else {
		return v, nil
	}
}

func (pwc *Pool_conn_obj) Close() error {
	pwc.pc.put(pwc)
	return nil
}
