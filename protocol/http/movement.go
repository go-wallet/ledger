package http

import (
	"time"

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

type GetMovementResponse struct {
	ID               string `json:"id"`
	AccountID        string `json:"account_id"`
	IsDebit          bool   `json:"is_debit"`
	Amount           int    `json:"amount"`
	PreviousMovement string `json:"previous_movement"`
	PreviousBalance  int    `json:"previous_balance"`
	CreatedAt        string `json:"created_at"`
}

type GetMovementsResponse struct {
	Data []*GetMovementResponse `json:"data"`
}

func ResponseFromMovement(mov *movement.Movement) *GetMovementResponse {
	return &GetMovementResponse{
		ID:               mov.ID.String(),
		AccountID:        mov.AccountID.String(),
		IsDebit:          mov.IsDebit,
		Amount:           mov.Amount,
		PreviousMovement: mov.PreviousMovement.String(),
		PreviousBalance:  mov.PreviousBalance,
		CreatedAt:        mov.CreatedAt.Format(time.RFC3339),
	}
}
