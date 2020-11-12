package movement

import "fmt"

type AlreadyExistsError struct {
	id ID
}

func NewAlreadyExistsError(id ID) *AlreadyExistsError {
	return &AlreadyExistsError{
		id: id,
	}
}

func (err *AlreadyExistsError) Error() string {
	return fmt.Sprintf("movement already exists: %s", err.id.String())
}
