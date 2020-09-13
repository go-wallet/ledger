package account

import (
	"context"
	"errors"
	"math"
	"sync"
	"time"
)

const Retries = 3
const WaitFor = 1 * time.Second

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
	retries int
	waitFor time.Duration
}

func NewLockerClient(l Lockable) *LockerClient {
	return &LockerClient{
		locker: l,
	}
}

func NewLocker(clis []*LockerClient) *Locker {
	l := &Locker{clients: clis}

	return l.WithRetries(Retries).WithWaitFor(WaitFor)
}

func (l *Locker) WithRetries(retries int) *Locker {
	return &Locker{
		clients: l.clients,
		retries: retries,
		waitFor: l.waitFor,
	}
}

func (l *Locker) WithWaitFor(waitFor time.Duration) *Locker {
	return &Locker{
		clients: l.clients,
		retries: l.retries,
		waitFor: waitFor,
	}
}

func (l *Locker) Lock(ctx context.Context, id ID, key string) error {
	var err error
	for c := 1; c < l.retries; c++ {
		err = l.doLock(ctx, id, key)
		if err == nil {
			return nil
		}

		time.Sleep(l.waitFor)
	}

	return err
}

func (l *Locker) Unlock(ctx context.Context, id ID, key string) error {
	buf := len(l.clients)
	wg := sync.WaitGroup{}
	wg.Add(buf)

	unlockFn := func(l Lockable) {
		defer wg.Done()
		l.Unlock(ctx, id, key)
	}

	for i := range l.clients {
		go unlockFn(l.clients[i].locker)
	}

	wg.Wait()

	return nil
}

func (l *Locker) doLock(ctx context.Context, id ID, key string) error {
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
