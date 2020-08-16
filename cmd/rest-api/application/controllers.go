package application

import (
	"net/http"
	"time"

	"github.com/labstack/echo"

	"github.com/vsmoraes/open-ledger/ledger"
	"github.com/vsmoraes/open-ledger/ledger/account"
	"github.com/vsmoraes/open-ledger/ledger/movement"
)

type FindMovementsRequest struct {
	AccountID string `query:"account_id"`
}

type CreateMovementRequest struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	IsDebit   bool   `json:"is_debit"`
	Amount    int    `json:"amount"`
}

func findMomentsController(mf movement.FindableByAccount) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		req := &FindMovementsRequest{}
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

		return ctx.JSON(http.StatusOK, movements)
	}
}

func createMovementController(ledger *ledger.Ledger) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		req := &CreateMovementRequest{}
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

		if err := ledger.AddEntry(ctx.Request().Context(), entry); err != nil {
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}

		return ctx.JSON(http.StatusCreated, nil)
	}
}
