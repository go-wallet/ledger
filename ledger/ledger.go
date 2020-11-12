package ledger

import (
	"context"

	"github.com/vsmoraes/open-ledger/ledger/account"
	"github.com/vsmoraes/open-ledger/ledger/movement"
)

type Ledger struct {
	accountLocker   *account.Locker
	movementFinder  movement.FindableByAccount
	movementCreator movement.Creatable
	movmentFetcher  movement.Fetchable
}

func New(al *account.Locker, tf movement.FindableByAccount, tc movement.Creatable, mf movement.Fetchable) *Ledger {
	return &Ledger{
		accountLocker:   al,
		movementFinder:  tf,
		movementCreator: tc,
		movmentFetcher:  mf,
	}
}

func (l *Ledger) CreateMovement(ctx context.Context, m movement.Movement) error {
	defer l.accountLocker.Unlock(ctx, m.AccountID, m.ID.String())

	if conflict, err := l.movmentFetcher.Fetch(ctx, m.ID); err == nil && conflict != nil {
		return movement.NewAlreadyExistsError(m.ID)
	}

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
		return NewBalanceNotEnoughError(lastEntry.CurrentBalance(), newEntry.Amount)
	}

	return nil
}
