package application

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/vsmoraes/open-ledger/config"
	protocol "github.com/vsmoraes/open-ledger/protocol/http"
)

type ControllersTestSuite struct {
	suite.Suite

	app *Application

	movToCreate *protocol.CreateMovementRequest
}

func (cts *ControllersTestSuite) SetupSuite() {
	cts.app = NewApplication()
	cts.movToCreate = &protocol.CreateMovementRequest{
		ID:        uuid.NewV4().String(),
		AccountID: uuid.NewV4().String(),
		IsDebit:   false,
		Amount:    100000,
	}

	go cts.app.Start(config.Config().RestAPI.Port)
}

func (cts *ControllersTestSuite) TearDownSuite() {
	cts.app.stop()
}

func (cts *ControllersTestSuite) TestCreateSingleMovement() {
	body, _ := json.Marshal(cts.movToCreate)
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(fmt.Sprintf("http://localhost%s/ledger/movements", config.Config().RestAPI.Port))

	assert.Nil(cts.T(), err)
	assert.Equal(cts.T(), http.StatusCreated, resp.StatusCode())
}

func (cts *ControllersTestSuite) TestGetCreatedMovement() {
	result := &protocol.GetMovementsResponse{
		Data: make([]*protocol.GetMovementResponse, 0),
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetResult(result).
		Get(fmt.Sprintf("http://localhost%s/ledger/movements?account_id=%s", config.Config().RestAPI.Port, cts.movToCreate.AccountID))

	assert.Nil(cts.T(), err)
	assert.Equal(cts.T(), http.StatusOK, resp.StatusCode())
	assert.Len(cts.T(), result.Data, 1)
	assert.Equal(cts.T(), cts.movToCreate.ID, result.Data[0].ID)
	assert.Equal(cts.T(), cts.movToCreate.AccountID, result.Data[0].AccountID)
	assert.Equal(cts.T(), cts.movToCreate.IsDebit, result.Data[0].IsDebit)
	assert.Equal(cts.T(), cts.movToCreate.Amount, result.Data[0].Amount)
}

func TestControllersTestSuite(t *testing.T) {
	suite.Run(t, &ControllersTestSuite{})
}
