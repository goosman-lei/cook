package pool

import (
	"crypto/tls"
	"time"
)

type Pool_tls_conn struct {
	*Pool
}

type Pool_tls_conn_obj struct {
	tls.Conn
	ptc *Pool_tls_conn
}

func NewPool_tls_conn(c int, factory func() (tls.Conn, error)) *Pool_tls_conn {
	ptc := &Pool_tls_conn{}
	ptc.Pool = new_pool(c, func() (interface{}, error) {
		if c, err := factory(); err != nil {
			return nil, err
		} else {
			return &Pool_tls_conn_obj{
				Conn: c,
				ptc:  ptc,
			}, nil
		}
	}, func(d interface{}) {
		if c, ok := d.(*Pool_tls_conn_obj); ok {
			c.Conn.Close()
		}
	})

	return ptc
}

func (ptc *Pool_tls_conn) Get() (*Pool_tls_conn_obj, error) {
	if c, err := ptc.get(); err != nil {
		return nil, err
	} else if v, ok := c.(*Pool_tls_conn_obj); !ok {
		return nil, ErrNotWrapper
	} else {
		return v, nil
	}
}

func (ptc *Pool_tls_conn) Get_timeout(duration time.Duration) (*Pool_tls_conn_obj, error) {
	if c, err := ptc.get_timeout(duration); err != nil {
		return nil, err
	} else if v, ok := c.(*Pool_tls_conn_obj); !ok {
		return nil, ErrNotWrapper
	} else {
		return v, nil
	}
}

func (pwc *Pool_tls_conn_obj) Close() error {
	pwc.ptc.put(pwc)
	return nil
}
