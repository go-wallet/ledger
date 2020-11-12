package application

type BadRequestError struct {
	message string
}

type InvalidRequestDataError struct {
	message string
}

func NewBadRequestError(msg string) *BadRequestError {
	return &BadRequestError{message: msg}
}

func NewInvalidRequestDataError(msg string) *InvalidRequestDataError {
	return &InvalidRequestDataError{message: msg}
}

func (err *BadRequestError) Error() string {
	return err.message
}

func (err *InvalidRequestDataError) Error() string {
	return err.message
}
