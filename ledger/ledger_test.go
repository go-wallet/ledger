package ledger

import (
	"context"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"syreclabs.com/go/faker"

	"github.com/vsmoraes/open-ledger/ledger/account"
	"github.com/vsmoraes/open-ledger/ledger/mock"
	"github.com/vsmoraes/open-ledger/ledger/movement"
)

type LedgerTestSuite struct {
	suite.Suite
}

func (lts *LedgerTestSuite) TestCreateNewEntry() {
	accountID := account.ID(uuid.NewV4().String())
	m := movement.Movement{
		ID:               movement.ID(uuid.NewV4().String()),
		AccountID:        accountID,
		IsDebit:          false,
		Amount:           faker.Number().NumberInt(3),
		PreviousMovement: "",
		PreviousBalance:  0,
		CreatedAt:        time.Now(),
	}

	al := &mock.AccountLocker{}
	al.On("Lock", accountID).Return(nil)
	al.On("Unlock", accountID).Return(nil)

	mf := &mock.MovementFinder{}
	mf.On("Last", accountID).Return(nil, nil)

	mc := &mock.MovementCreator{}
	mc.On("Create", &m).Return(nil)

	l := New(al, mf, mc)

	result := l.AddEntry(context.Background(), m)
	assert.Nil(lts.T(), result)
	al.AssertExpectations(lts.T())
	mf.AssertExpectations(lts.T())
	mc.AssertExpectations(lts.T())
}

func (lts *LedgerTestSuite) TestFailWhenValidationDoesNotPass()           {}
func (lts *LedgerTestSuite) TestFailWhenLockingIsNotPossible()            {}
func (lts *LedgerTestSuite) TestFailWhenErrorIsReceivedFromLastMovement() {}
func (lts *LedgerTestSuite) TestFailWhenErrorIsReceivedFromCreation()     {}

func TestLedgerTestSuite(t *testing.T) {
	suite.Run(t, &LedgerTestSuite{})
}
