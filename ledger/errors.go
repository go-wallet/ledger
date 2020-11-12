package ledger

import "fmt"

type BalanceNotEnoughError struct {
	currBalance   int
	neededBalance int
}

func NewBalanceNotEnoughError(balance, needed int) *BalanceNotEnoughError {
	return &BalanceNotEnoughError{
		currBalance:   balance,
		neededBalance: needed,
	}
}

func (err *BalanceNotEnoughError) Error() string {
	return fmt.Sprintf("current balance (%d) is not enough to complete this operation (%d)", err.currBalance, err.neededBalance)
}
