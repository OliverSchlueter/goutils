package broker

import "github.com/nats-io/nats.go"

type Broker interface {
	Publish(subject string, data []byte) error
	Request(subject string, data []byte) (*nats.Msg, error)
	Subscribe(subject string, handler nats.MsgHandler) error
	SubscribeQueue(subject, queue string, handler nats.MsgHandler) error
}
