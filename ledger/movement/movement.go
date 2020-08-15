package movement

import (
	"context"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/vsmoraes/open-ledger/ledger/account"
)

type FindableByAccount interface {
	All(ctx context.Context, id account.ID) ([]*Movement, error)
	Last(ctx context.Context, id account.ID) (*Movement, error)
}

type Creatable interface {
	Create(ctx context.Context, t *Movement) error
}

type Movement struct {
	ID              uuid.UUID
	Account         *account.Account
	IsDebit         bool
	Amount          int
	PreviousBalance int
	CreatedAt       time.Time
}

func (m Movement) CurrentBalance() int {
	if m.IsDebit {
		return m.PreviousBalance - m.Amount
	}

	return m.PreviousBalance + m.Amount
}
