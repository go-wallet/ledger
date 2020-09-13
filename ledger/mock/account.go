package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/vsmoraes/open-ledger/ledger/account"
)

type LockerClient struct {
	mock.Mock
}

func (alm *LockerClient) Lock(ctx context.Context, id account.ID, key string) error {
	args := alm.Called(id)

	return args.Error(0)
}

func (alm *LockerClient) Unlock(ctx context.Context, id account.ID, key string) error {
	args := alm.Called(id)

	return args.Error(0)
}
