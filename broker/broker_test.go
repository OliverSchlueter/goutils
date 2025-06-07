package broker_test

import (
	"context"
	"github.com/OliverSchlueter/goutils/broker"
	"github.com/OliverSchlueter/goutils/broker/brokertest"
	"github.com/OliverSchlueter/goutils/containers"
	"github.com/nats-io/nats.go"
	"testing"
)

var natsBroker broker.Broker

func TestFakeBroker_Publish(t *testing.T) {
	brokertest.TestPublish(t, broker.NewFakeBroker())
}

func TestFakeBroker_Request(t *testing.T) {
	brokertest.TestRequest(t, broker.NewFakeBroker())
}

func TestFakeBroker_Subscribe(t *testing.T) {
	brokertest.TestSubscribe(t, broker.NewFakeBroker())
}

func TestFakeBroker_SubscribeQueue(t *testing.T) {
	brokertest.TestSubscribeQueue(t, broker.NewFakeBroker())
}

func TestNatsBroker_Publish(t *testing.T) {
	brokertest.TestPublish(t, NewNatsBroker(t))
}

func TestNatsBroker_Request(t *testing.T) {
	brokertest.TestRequest(t, NewNatsBroker(t))
}

func TestNatsBroker_Subscribe(t *testing.T) {
	brokertest.TestSubscribe(t, NewNatsBroker(t))
}

func TestNatsBroker_SubscribeQueue(t *testing.T) {
	brokertest.TestSubscribeQueue(t, NewNatsBroker(t))
}

func NewNatsBroker(t *testing.T) broker.Broker {
	if natsBroker != nil {
		return natsBroker
	}

	ctx := context.Background()
	_, err := containers.StartNATS(ctx)
	if err != nil {
		t.Fatalf("Failed to start NATS container: %v", err)
	}

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Fatalf("Failed to connect to NATS: %v", err)
	}

	natsBroker = broker.NewNatsBroker(&broker.NatsConfiguration{
		Nats: nc,
	})

	return natsBroker
}
