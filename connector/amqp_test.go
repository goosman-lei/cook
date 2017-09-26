package connector

import (
	"github.com/streadway/amqp"
	"testing"
)

func TestCase(t *testing.T) {
	var (
		err  error
		conn *AMQPWrapper
		ch   <-chan amqp.Delivery
	)

	SetupAmqp(map[string]AMQPConf{
		"default": AMQPConf{Username: "nice", Password: "zp50MpILYUono", Addr: "rabbitmq.niceprivate.com:5672", Vhost: "/"},
	})
	conn = MustGetAmqp("default")
	conn.Produce("show-pub-1", "guoguo.amqp.exchange", "RK_1")
	conn.Produce("show-pub-2", "guoguo.amqp.exchange", "RK_1")
	conn.Produce("show-pub-3", "guoguo.amqp.exchange", "RK_1")
	conn.Produce("show-pub-4", "guoguo.amqp.exchange", "RK_1")
	conn.Produce("show-pub-5", "guoguo.amqp.exchange", "RK_1")

	if ch, err = conn.Consume("guoguo.amqp.queue"); err != nil {
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
