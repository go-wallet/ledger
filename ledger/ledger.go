package ledger

import (
	"context"
	"errors"

	"github.com/vsmoraes/open-ledger/ledger/account"
	"github.com/vsmoraes/open-ledger/ledger/movement"
)

var ErrBalanceNotEnough = errors.New("current balance is not enough to complete this operation")

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

func (l *Ledger) CreateMovement(ctx context.Context, m movement.Movement) error {
	defer l.accountLocker.Unlock(ctx, m.AccountID, m.ID.String())

	if err := l.accountLocker.Lock(ctx, m.AccountID, m.ID.String()); err != nil {
		return err
	}

	lst, err := l.movementFinder.Last(ctx, m.AccountID)
	if err != nil {
		return err
	}

	if err = l.validateEntry(&m, lst); err != nil {
		return err
	}

	if lst != nil {
		m.PreviousMovement = lst.ID
		m.PreviousBalance = lst.CurrentBalance()
	}

	if err = l.movementCreator.Create(ctx, &m); err != nil {
		return err
	}

	return nil
}

func (l *Ledger) validateEntry(newEntry, lastEntry *movement.Movement) error {
	if newEntry.IsDebit && (lastEntry == nil || lastEntry.CurrentBalance() < newEntry.Amount) {
		return ErrBalanceNotEnough
	}

	return nil
}
