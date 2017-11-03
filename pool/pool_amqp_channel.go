package pool

import (
	"github.com/streadway/amqp"
	"time"
)

type Pool_amqp_channel struct {
	*Pool
}

type Pool_amqp_channel_obj struct {
	*amqp.Channel
	pac *Pool_amqp_channel
}

func NewPool_amqp_channel(c int, factory func() (*amqp.Channel, error)) *Pool_amqp_channel {
	pac := &Pool_amqp_channel{}
	pac.Pool = new_pool(c, func() (interface{}, error) {
		if c, err := factory(); err != nil {
			return nil, err
		} else {
			return &Pool_amqp_channel_obj{
				Channel: c,
				pac:     pac,
			}, nil
		}
	}, func(d interface{}) {
		if c, ok := d.(*Pool_amqp_channel_obj); ok {
			c.Channel.Close()
		}
	})

	return pac
}

func (pac *Pool_amqp_channel) Get() (*Pool_amqp_channel_obj, error) {
	if c, err := pac.get(); err != nil {
		return nil, err
	} else if v, ok := c.(*Pool_amqp_channel_obj); !ok {
		return nil, ErrNotWrapper
	} else {
		return v, nil
	}
}

func (pac *Pool_amqp_channel) Get_timeout(duration time.Duration) (*Pool_amqp_channel_obj, error) {
	if c, err := pac.get_timeout(duration); err != nil {
		return nil, err
	} else if v, ok := c.(*Pool_amqp_channel_obj); !ok {
		return nil, ErrNotWrapper
	} else {
		return v, nil
	}
}

func (p *Pool_amqp_channel_obj) Close() error {
	p.pac.put(p)
	return nil
}

func (p *Pool_amqp_channel_obj) Produce_normal(payload, exchange, routingKey string) error {
	return p.Channel.Publish(exchange, routingKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(payload),
	})
}

func (p *Pool_amqp_channel_obj) Consume_normal(queue string) (<-chan amqp.Delivery, error) {
	return p.Consume(queue, "", false, false, false, false, nil)
}
