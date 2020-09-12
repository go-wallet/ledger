package account

import (
	"context"
)

type ID string

func (id ID) String() string {
	return string(id)
}

type Lockable interface {
	Lock(ctx context.Context, id ID, key string) error
	Unlock(ctx context.Context, id ID, key string) error
}
