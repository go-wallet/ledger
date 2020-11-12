package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/vsmoraes/open-ledger/ledger/account"
	"github.com/vsmoraes/open-ledger/ledger/movement"
)

type MovementFinder struct {
	mock.Mock
}

func (mf *MovementFinder) All(ctx context.Context, id account.ID) ([]*movement.Movement, error) {
	args := mf.Called(id)

	return args.Get(0).([]*movement.Movement), args.Error(1)
}

func (mf *MovementFinder) Last(ctx context.Context, id account.ID) (*movement.Movement, error) {
	args := mf.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*movement.Movement), args.Error(1)
}

type MovementCreator struct {
	mock.Mock
}

func (mc *MovementCreator) Create(ctx context.Context, m *movement.Movement) error {
	args := mc.Called(m)

	return args.Error(0)
}

type MovementFetcher struct {
	mock.Mock
}

func (mf *MovementFetcher) Fetch(ctx context.Context, id movement.ID) (*movement.Movement, error) {
	args := mf.Called(id)

	return args.Get(0).(*movement.Movement), args.Error(1)
}
