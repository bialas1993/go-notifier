package event

import (
	"bytes"
	"encoding/gob"

	"github.com/nats-io/go-nats"
	"github.com/bialas1993/go-notifier/schema"
)

type NatsEventStore struct {
	nc                        *nats.Conn
	notifyCreatedSubscription *nats.Subscription
	notifyCreatedChan         chan NotifyCreatedMessage
}

func NewNats(url string) (*NatsEventStore, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsEventStore{nc: nc}, nil
}

func (e *NatsEventStore) SubscribeNotifyCreated() (<-chan NotifyCreatedMessage, error) {
	m := NotifyCreatedMessage{}
	e.notifyCreatedChan = make(chan NotifyCreatedMessage, 64)
	ch := make(chan *nats.Msg, 64)
	var err error
	e.notifyCreatedSubscription, err = e.nc.ChanSubscribe(m.Key(), ch)
	if err != nil {
		return nil, err
	}
	// Decode message
	go func() {
		for {
			select {
			case msg := <-ch:
				e.readMessage(msg.Data, &m)
				e.notifyCreatedChan <- m
			}
		}
	}()
	return (<-chan NotifyCreatedMessage)(e.notifyCreatedChan), nil
}

func (e *NatsEventStore) OnNotifyCreated(f func(NotifyCreatedMessage)) (err error) {
	m := NotifyCreatedMessage{}
	e.notifyCreatedSubscription, err = e.nc.Subscribe(m.Key(), func(msg *nats.Msg) {
		e.readMessage(msg.Data, &m)
		f(m)
	})
	return
}

func (e *NatsEventStore) Close() {
	if e.nc != nil {
		e.nc.Close()
	}
	if e.notifyCreatedSubscription != nil {
		e.notifyCreatedSubscription.Unsubscribe()
	}
	close(e.notifyCreatedChan)
}

func (e *NatsEventStore) PublishNotifyCreated(notify schema.Notify) error {
	m := NotifyCreatedMessage{
		notify.ID,
		notify.Title,
		notify.Body,
		notify.Service,
		notify.CreatedAt,
	}
	data, err := e.writeMessage(&m)
	if err != nil {
		return err
	}
	return e.nc.Publish(m.Key(), data)
}

func (mq *NatsEventStore) writeMessage(m Message) ([]byte, error) {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(m)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (mq *NatsEventStore) readMessage(data []byte, m interface{}) error {
	b := bytes.Buffer{}
	b.Write(data)
	return gob.NewDecoder(&b).Decode(m)
}
