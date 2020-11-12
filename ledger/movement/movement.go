package movement

import (
	"context"
	"time"

	"github.com/vsmoraes/open-ledger/ledger/account"
)

type Fetchable interface {
	Fetch(ctx context.Context, id ID) (*Movement, error)
}

type FindableByAccount interface {
	All(ctx context.Context, id account.ID) ([]*Movement, error)
	Last(ctx context.Context, id account.ID) (*Movement, error)
}

type Creatable interface {
	Create(ctx context.Context, m *Movement) error
}

type ID string

func (id ID) String() string {
	return string(id)
}

type Movement struct {
	ID               ID
	AccountID        account.ID
	IsDebit          bool
	Amount           int
	PreviousMovement ID
	PreviousBalance  int
	CreatedAt        time.Time
}

func (m Movement) CurrentBalance() int {
	if m.IsDebit {
		return m.PreviousBalance - m.Amount
	}

	return m.PreviousBalance + m.Amount
}
