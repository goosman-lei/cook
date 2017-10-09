package connector

import (
	"fmt"
	"github.com/streadway/amqp"
	cook_log "gitlab.niceprivate.com/golang/cook/log"
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
}

type AMQPWrapper struct {
	Url  string
	Conn *amqp.Connection
}

var (
	amqpConnMapping *cook_util.CMap
)

func SetupAmqp(configs map[string]AMQPConf) error {
	var (
		url  string
		err  error
		conn *amqp.Connection
	)
	amqpConnMapping = cook_util.NewCMap()
	for sn, config := range configs {
		url = fmt.Sprintf("amqp://%s:%s@%s%s", config.Username, config.Password, config.Addr, config.Vhost)
		if conn, err = amqp.DialConfig(url, amqp.Config{
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, config.ConnectTimeout)
			},
		}); err != nil {
			return err
		}
		amqpConnMapping.Set(sn, &AMQPWrapper{
			Url:  url,
			Conn: conn,
		})
	}
	return nil
}

func GetAmqp(sn string) (*AMQPWrapper, error) {
	if conn, exists := amqpConnMapping.Get(sn); exists {
		return conn.(*AMQPWrapper), nil
	}
	cook_log.Warnf("get amqp [%s], but not ready", sn)
	return nil, fmt.Errorf("have no amqp: %s", sn)
}

func MustGetAmqp(sn string) *AMQPWrapper {
	conn, err := GetAmqp(sn)
	if err != nil {
		panic(err)
	}
	return conn
}

func (c *AMQPWrapper) Produce(payload, exchange, routingKey string) error {
	if channel, err := c.Conn.Channel(); err == nil {
		return channel.Publish(exchange, routingKey, false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(payload),
		})
	} else {
		return err
	}
}

func (c *AMQPWrapper) Consume(queue string) (<-chan amqp.Delivery, error) {
	if channel, err := c.Conn.Channel(); err == nil {
		return channel.Consume(queue, "", true, false, false, false, nil)
	} else {
		return nil, err
	}
}
