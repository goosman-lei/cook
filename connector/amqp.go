package connector

import (
	"fmt"
	"github.com/streadway/amqp"
	cook_pool "gitlab.niceprivate.com/golang/cook/pool"
	cook_util "gitlab.niceprivate.com/golang/cook/util"
	"net"
	"time"
)

type AMQPConf struct {
	Username string
	Password string
	Addr     string
	Vhost    string

	ConnectTimeout time.Duration

	MaxConn    int
	MaxChannel int // must greater or equal to MaxConn
}

var (
	amqp_pool_mapping *cook_util.CMap = cook_util.NewCMap()
)

func SetupAmqp(configs map[string]AMQPConf) error {
	for sn, config := range configs {
		if config.MaxChannel < config.MaxConn {
			config.MaxChannel = config.MaxConn
		}

		conn_pool := func(c AMQPConf) *cook_pool.Pool_amqp_conn {
			url := fmt.Sprintf("amqp://%s:%s@%s/%s", c.Username, c.Password, c.Addr, c.Vhost)
			return cook_pool.NewPool_amqp_conn(c.MaxConn, func() (*amqp.Connection, error) {
				return amqp.DialConfig(
					url,
					amqp.Config{
						Dial: func(network, addr string) (net.Conn, error) {
							return net.DialTimeout(network, addr, c.ConnectTimeout)
						},
					})
			})
		}(config)

		func(p *cook_pool.Pool_amqp_conn, c AMQPConf, s string) {
			amqp_pool_mapping.Set(s, cook_pool.NewPool_amqp_channel(c.MaxChannel, func() (*amqp.Channel, error) {
				if conn, err := p.Get(); err != nil {
					return nil, err
				} else {
					// must put back connection
					defer conn.Close()
					return conn.Channel()
				}
			}))
		}(conn_pool, config, sn)
	}
	return nil
}

func Get_amqp_obj(sn string) (*cook_pool.Pool_amqp_channel_obj, error) {
	if wrapper, exists := amqp_pool_mapping.Get(sn); !exists {
		return nil, fmt.Errorf("have no amqp: %s", sn)
	} else {
		return wrapper.(*cook_pool.Pool_amqp_channel).Get()
	}
}
