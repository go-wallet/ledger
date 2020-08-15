package ledger

import (
	"context"
	"errors"

	"github.com/vsmoraes/open-ledger/ledger/account"
	"github.com/vsmoraes/open-ledger/ledger/movement"
)

type Ledger struct {
	accountLocker   account.Lockable
	movementFinder  movement.FindableByAccount
	movementCreator movement.Creatable
}

func New(al account.Lockable, tf movement.FindableByAccount, tc movement.Creatable) *Ledger {
	return &Ledger{
		accountLocker:   al,
		movementFinder:  tf,
		movementCreator: tc,
	}
}

func (l *Ledger) AddEntry(ctx context.Context, m movement.Movement) error {
	defer l.accountLocker.Unlock(ctx, m.Account, m.ID.String())

	if err := l.accountLocker.Lock(ctx, m.Account, m.ID.String()); err != nil {
		return err
	}

	lst, err := l.movementFinder.Last(ctx, m.Account.ID)
	if err != nil {
		return err
	}

	if m.IsDebit && lst.CurrentBalance() < m.Amount {
		return errors.New("current balance is not available to complete this operation")
	}

	if err = l.movementCreator.Create(ctx, &m); err != nil {
		return err
	}

	return nil
}
