package account

import "fmt"

type NotEnoughQuorumError struct {
	quorum int
	needed int
}

func NewNotEnoughQuorumError(quorum, needed int) *NotEnoughQuorumError {
	return &NotEnoughQuorumError{
		quorum: quorum,
		needed: needed,
	}
}

func (err *NotEnoughQuorumError) Error() string {
	return fmt.Sprintf("not enough quorum (%d) for locking account, needed: %d", err.quorum, err.needed)
}
