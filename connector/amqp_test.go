package connector

import (
	"github.com/streadway/amqp"
	cook_pool "gitlab.niceprivate.com/golang/cook/pool"
	"testing"
)

func TestCase(t *testing.T) {
	var (
		err error
		mq  *cook_pool.Pool_amqp_channel_obj
		ch  <-chan amqp.Delivery
	)

	SetupAmqp(map[string]AMQPConf{
		"default": AMQPConf{Username: "nice", Password: "zp50MpILYUono", Addr: "rabbitmq.niceprivate.com:5672", Vhost: "/", MaxConn: 1, MaxChannel: 1},
	})
	mq, err = Get_amqp_obj("default")
	mq.Produce_normal("show-pub-1", "guoguo.amqp.exchange.fanout", "RK_1")
	mq.Produce_normal("show-pub-2", "guoguo.amqp.exchange.fanout", "RK_1")
	mq.Produce_normal("show-pub-3", "guoguo.amqp.exchange.fanout", "RK_1")
	mq.Produce_normal("show-pub-4", "guoguo.amqp.exchange.fanout", "RK_1")
	mq.Produce_normal("show-pub-5", "guoguo.amqp.exchange.fanout", "RK_1")

	if ch, err = mq.Consume_normal("guoguo.amqp.queue"); err != nil {
		t.Logf("consume failed: %s", err)
	}

	if string((<-ch).Body) != "show-pub-1" {
		t.Fatalf("msg 1 is not show-pub-1")
	}
	if string((<-ch).Body) != "show-pub-2" {
		t.Fatalf("msg 1 is not show-pub-2")
	}
	if string((<-ch).Body) != "show-pub-3" {
		t.Fatalf("msg 1 is not show-pub-3")
	}
	if string((<-ch).Body) != "show-pub-4" {
		t.Fatalf("msg 1 is not show-pub-4")
	}
	if string((<-ch).Body) != "show-pub-5" {
		t.Fatalf("msg 1 is not show-pub-5")
	}
}
