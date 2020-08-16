package account

import (
	"context"
)

type ID string

type Lockable interface {
	Lock(ctx context.Context, id ID, key string) error
	Unlock(ctx context.Context, id ID, key string) error
}

type Fetcher interface {
	Fetch(id ID) (ID, error)
}
