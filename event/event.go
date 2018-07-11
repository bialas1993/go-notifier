package event

import "github.com/bialas1993/go-notifier/schema"

type EventStore interface {
	Close()
	PublishNotifyCreated(notify schema.Notify) error
	SubscribeNotifyCreated() (<-chan NotifyCreatedMessage, error)
	OnNotifyCreated(f func(NotifyCreatedMessage)) error
}

var impl EventStore

func SetEventStore(es EventStore) {
	impl = es
}

func Close() {
	impl.Close()
}

func PublishNotifyCreated(notify schema.Notify) error {
	return impl.PublishNotifyCreated(notify)
}

func SubscribeNotifyCreated() (<-chan NotifyCreatedMessage, error) {
	return impl.SubscribeNotifyCreated()
}

func OnNotifyCreated(f func(NotifyCreatedMessage)) error {
	return impl.OnNotifyCreated(f)
}
