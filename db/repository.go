package db

import (
	"context"

	"github.com/bialas1993/go-notifier/schema"
)

type Repository interface {
	Close()
	Insert(ctx context.Context, notify schema.Notify) error
	List(ctx context.Context, skip uint64, take uint64) ([]schema.Notify, error)
}

var impl Repository

func SetRepository(repository Repository) {
	impl = repository
}

func Close() {
	impl.Close()
}

func Insert(ctx context.Context, notify schema.Notify) error {
	return impl.Insert(ctx, notify)
}

func List(ctx context.Context, skip uint64, take uint64) ([]schema.Notify, error) {
	return impl.List(ctx, skip, take)
}
