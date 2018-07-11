package search

import (
	"context"

	"github.com/bialas1993/go-notifier/schema"
)

type Repository interface {
	Close()
	Insert(ctx context.Context, notify schema.Notify) error
	Search(ctx context.Context, query string, skip uint64, take uint64) ([]schema.Notify, error)
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

func Search(ctx context.Context, query string, skip uint64, take uint64) ([]schema.Notify, error) {
	return impl.Search(ctx, query, skip, take)
}
