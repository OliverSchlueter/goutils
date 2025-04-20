package broker

import "github.com/nats-io/nats.go"

type FakeBroker struct {
	subscribers []func(msg *nats.Msg)
}

func NewFakeBroker() *FakeBroker {
	return &FakeBroker{}
}

func (b *FakeBroker) Publish(subject string, data []byte) error {
	msg := &nats.Msg{
		Subject: subject,
		Reply:   "reply",
		Header:  nats.Header{},
		Data:    data,
	}

	for _, s := range b.subscribers {
		s(msg)
	}

	return nil
}

func (b *FakeBroker) Subscribe(subject string, handler nats.MsgHandler) error {
	b.subscribers = append(b.subscribers, handler)
	return nil
}

func (b *FakeBroker) SubscribeQueue(subject, queue string, handler nats.MsgHandler) error {
	b.subscribers = append(b.subscribers, handler)
	return nil
}
