package broker

import (
	"fmt"
	"github.com/nats-io/nats.go"
)

type FakeBroker struct {
	subscribers    []func(msg *nats.Msg)
	requestHandler RequestHandler
}

type RequestHandler func(msg *nats.Msg) (*nats.Msg, error)

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

func (b *FakeBroker) Request(subject string, data []byte) (*nats.Msg, error) {
	msg := &nats.Msg{
		Subject: subject,
		Reply:   "reply",
		Header:  nats.Header{},
		Data:    data,
	}

	if b.requestHandler != nil {
		return b.requestHandler(msg)
	}

	return nil, fmt.Errorf("no request handler set")
}

func (b *FakeBroker) SetRequestHandler(handler RequestHandler) {
	b.requestHandler = handler
}

func (b *FakeBroker) Subscribe(subject string, handler nats.MsgHandler) error {
	b.subscribers = append(b.subscribers, handler)
	return nil
}

func (b *FakeBroker) SubscribeQueue(subject, queue string, handler nats.MsgHandler) error {
	b.subscribers = append(b.subscribers, handler)
	return nil
}
