package account

import (
	"context"
	"errors"
	"math"
	"sync"
)

var ErrNotEnoughQuorum = errors.New("not enough quorum for locking")

type Lockable interface {
	Lock(ctx context.Context, id ID, key string) error
	Unlock(ctx context.Context, id ID, key string) error
}

type LockerClient struct {
	locker Lockable
}

type Locker struct {
	clients []*LockerClient
}

func NewLockerClient(l Lockable) *LockerClient {
	return &LockerClient{
		locker: l,
	}
}

func NewLocker(clis []*LockerClient) *Locker {
	return &Locker{clients: clis}
}

func (l *Locker) Lock(ctx context.Context, id ID, key string) error {
	buf := len(l.clients)
	c := float64(buf)
	quorum := int(math.Min(c, c/2+1))

	locks := make(chan bool, buf)
	wg := sync.WaitGroup{}
	wg.Add(buf)

	lockFn := func(l Lockable) {
		defer wg.Done()
		if err := l.Lock(ctx, id, key); err == nil {
			locks <- true
		}
	}

	for i := range l.clients {
		go lockFn(l.clients[i].locker)
	}
	wg.Wait()

	if len(locks) >= quorum {
		return nil
	}

	close(locks)
	return ErrNotEnoughQuorum
}

func (l *Locker) Unlock(ctx context.Context, id ID, key string) error {
	for _, cli := range l.clients {
		cli.locker.Unlock(ctx, id, key)
	}

	return nil
}
