package account

import (
	"context"
)

type ID string

type Account struct {
	ID ID
}

type Lockable interface {
	Lock(ctx context.Context, account *Account, key string) error
	Unlock(ctx context.Context, account *Account, key string) error
}

type Fetcher interface {
	Fetch(id ID) (*Account, error)
}
