package broker

import "github.com/nats-io/nats.go"

type NatsBroker struct {
	nats *nats.Conn
}

type NatsConfiguration struct {
	Nats *nats.Conn
}

func NewNatsBroker(cfg *NatsConfiguration) *NatsBroker {
	return &NatsBroker{
		nats: cfg.Nats,
	}
}

func (b *NatsBroker) Publish(subject string, data []byte) error {
	return b.nats.Publish(subject, data)
}

func (b *NatsBroker) Request(subject string, data []byte) (*nats.Msg, error) {
	return b.nats.Request(subject, data, nats.DefaultTimeout)
}

func (b *NatsBroker) Subscribe(subject string, handler nats.MsgHandler) error {
	_, err := b.nats.Subscribe(subject, handler)
	return err
}

func (b *NatsBroker) SubscribeQueue(subject, queue string, handler nats.MsgHandler) error {
	_, err := b.nats.QueueSubscribe(subject, queue, handler)
	return err
}
