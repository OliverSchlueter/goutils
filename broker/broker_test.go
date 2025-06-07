package broker_test

import (
	"github.com/OliverSchlueter/goutils/broker"
	"github.com/OliverSchlueter/goutils/broker/brokertest"
	"testing"
)

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
