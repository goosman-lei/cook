package pool

import (
	"github.com/streadway/amqp"
	"time"
)

type Pool_amqp_conn struct {
	*Pool
}

type Pool_amqp_conn_obj struct {
	*amqp.Connection
	pac *Pool_amqp_conn
}

func NewPool_amqp_conn(c int, factory func() (*amqp.Connection, error)) *Pool_amqp_conn {
	pac := &Pool_amqp_conn{}
	pac.Pool = new_pool(c, func() (interface{}, error) {
		if c, err := factory(); err != nil {
			return nil, err
		} else {
			return &Pool_amqp_conn_obj{
				Connection: c,
				pac:        pac,
			}, nil
		}
	}, func(d interface{}) {
		if c, ok := d.(*Pool_amqp_conn_obj); ok {
			c.Connection.Close()
		}
	})

	return pac
}

func (pac *Pool_amqp_conn) Get() (*Pool_amqp_conn_obj, error) {
	if c, err := pac.get(); err != nil {
		return nil, err
	} else if v, ok := c.(*Pool_amqp_conn_obj); !ok {
		return nil, ErrNotWrapper
	} else {
		return v, nil
	}
}

func (pac *Pool_amqp_conn) Get_timeout(duration time.Duration) (*Pool_amqp_conn_obj, error) {
	if c, err := pac.get_timeout(duration); err != nil {
		return nil, err
	} else if v, ok := c.(*Pool_amqp_conn_obj); !ok {
		return nil, ErrNotWrapper
	} else {
		return v, nil
	}
}

func (p *Pool_amqp_conn_obj) Close() error {
	p.pac.put(p)
	return nil
}
