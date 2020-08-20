package application

import (
	"net/http"
	"time"

	"github.com/labstack/echo"

	"github.com/vsmoraes/open-ledger/cmd/rest-api/protocol"
	"github.com/vsmoraes/open-ledger/ledger"
	"github.com/vsmoraes/open-ledger/ledger/account"
	"github.com/vsmoraes/open-ledger/ledger/movement"
)

func findMomentsController(mf movement.FindableByAccount) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		req := &protocol.FindMovementsRequest{}
		if err := ctx.Bind(req); err != nil {
			return ctx.JSON(http.StatusBadRequest, err.Error())
		}
		if req.AccountID == "" {
			return ctx.JSON(http.StatusBadRequest, struct {
				Error string `json:"error"`
			}{
				Error: "You need to provide an account_id via query string",
			})
		}

		movements, err := mf.All(ctx.Request().Context(), account.ID(req.AccountID))
		if err != nil {
			return ctx.JSON(http.StatusNotFound, err.Error())
		}

		result := &protocol.GetMovementsResponse{
			Data: make([]*protocol.GetMovementResponse, 0),
		}
		for _, mov := range movements {
			result.Data = append(result.Data, protocol.ResponseFromMovement(mov))
		}
		return ctx.JSON(http.StatusOK, result)
	}
}

func createMovementController(ledger *ledger.Ledger) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		req := &protocol.CreateMovementRequest{}
		if err := ctx.Bind(req); err != nil {
			return ctx.JSON(http.StatusUnprocessableEntity, err.Error())
		}

		entry := movement.Movement{
			ID:        movement.ID(req.ID),
			AccountID: account.ID(req.AccountID),
			IsDebit:   req.IsDebit,
			Amount:    req.Amount,
			CreatedAt: time.Now(),
		}

		if err := ledger.CreateMovement(ctx.Request().Context(), entry); err != nil {
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}

		return ctx.JSON(http.StatusCreated, nil)
	}
}
