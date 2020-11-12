package application

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/vsmoraes/open-ledger/ledger/account"
	"github.com/vsmoraes/open-ledger/ledger/movement"

	"github.com/vsmoraes/open-ledger/ledger"
	protocol "github.com/vsmoraes/open-ledger/protocol/http"
)

func Error(ctx echo.Context, err error) error {
	var statusCode int
	var msg string

	switch t := err.(type) {
	case *ledger.BalanceNotEnoughError, *InvalidRequestDataError:
		statusCode = http.StatusUnprocessableEntity
		msg = t.Error()

	case *account.NotEnoughQuorumError:
		statusCode = http.StatusLocked
		msg = t.Error()

	case *movement.AlreadyExistsError:
		statusCode = http.StatusConflict
		msg = t.Error()

	case *BadRequestError:
		statusCode = http.StatusBadRequest
		msg = t.Error()

	default:
		statusCode = http.StatusInternalServerError
		msg = t.Error()
	}

	return ctx.JSON(statusCode, protocol.Error{Message: msg})
}
