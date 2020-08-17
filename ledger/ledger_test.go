package ledger

import (
	"context"
	"errors"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/vsmoraes/open-ledger/ledger/account"
	"github.com/vsmoraes/open-ledger/ledger/mock"
	"github.com/vsmoraes/open-ledger/ledger/movement"
)

type asserts func(t *testing.T, al *mock.AccountLocker, mf *mock.MovementFinder, mc *mock.MovementCreator, result error)

type LedgerTestSuite struct {
	suite.Suite
}

func (lts *LedgerTestSuite) TestCreateMovement() {
	tests := []struct {
		accountID    account.ID
		newMov       *movement.Movement
		lstMov       *movement.Movement
		lstError     error
		lockReturn   error
		unlockReturn error
		createError  error
		asserts      asserts
	}{
		// scenario 1: Create a new movement
		{
			accountID: "123-123-123",
			newMov: &movement.Movement{
				ID:        movement.ID("1234567890"),
				AccountID: "123-123-123",
				IsDebit:   false,
				Amount:    1000,
			},
			lstMov:       nil,
			lstError:     nil,
			lockReturn:   nil,
			unlockReturn: nil,
			createError:  nil,
			asserts: func(t *testing.T, al *mock.AccountLocker, mf *mock.MovementFinder, mc *mock.MovementCreator, result error) {
				assert.Nil(t, result)
				al.AssertExpectations(t)
				mf.AssertExpectations(t)
				mc.AssertExpectations(t)
			},
		},

		// scenario 2: Fail to create due to lack of balance
		{
			accountID: "123-123-123",
			newMov: &movement.Movement{
				ID:        movement.ID("1234567890"),
				AccountID: "123-123-123",
				IsDebit:   true,
				Amount:    1000,
			},
			lstMov: &movement.Movement{
				ID:               movement.ID(uuid.NewV4().String()),
				AccountID:        "123-123-123",
				IsDebit:          false,
				Amount:           500,
				PreviousMovement: "",
				PreviousBalance:  0,
				CreatedAt:        time.Now(),
			},
			lstError:     nil,
			lockReturn:   nil,
			unlockReturn: nil,
			createError:  nil,
			asserts: func(t *testing.T, al *mock.AccountLocker, mf *mock.MovementFinder, mc *mock.MovementCreator, result error) {
				assert.EqualError(t, result, ErrBalanceNotEnough.Error())
				al.AssertExpectations(t)
				mf.AssertExpectations(t)
				mc.AssertNumberOfCalls(t, "Create", 0)
			},
		},

		// scenario 3: Fail when creating a debit movement as the first one
		{
			accountID: "123-123-123",
			newMov: &movement.Movement{
				ID:        movement.ID("1234567890"),
				AccountID: "123-123-123",
				IsDebit:   true,
				Amount:    1000,
			},
			lstMov:       nil,
			lstError:     nil,
			lockReturn:   nil,
			unlockReturn: nil,
			createError:  nil,
			asserts: func(t *testing.T, al *mock.AccountLocker, mf *mock.MovementFinder, mc *mock.MovementCreator, result error) {
				assert.EqualError(t, result, ErrBalanceNotEnough.Error())
				al.AssertExpectations(t)
				mf.AssertExpectations(t)
				mc.AssertNumberOfCalls(t, "Create", 0)
			},
		},

		// scenario 4: Fail to create when locking is not possible
		{
			accountID: "123-123-123",
			newMov: &movement.Movement{
				ID:        movement.ID("1234567890"),
				AccountID: "123-123-123",
				IsDebit:   false,
				Amount:    1000,
			},
			lstMov:       nil,
			lstError:     nil,
			lockReturn:   errors.New("locking is not possible"),
			unlockReturn: nil,
			createError:  nil,
			asserts: func(t *testing.T, al *mock.AccountLocker, mf *mock.MovementFinder, mc *mock.MovementCreator, result error) {
				assert.EqualError(t, result, "locking is not possible")
				al.AssertExpectations(t)
				mf.AssertNumberOfCalls(t, "Last", 0)
				mc.AssertNumberOfCalls(t, "Create", 0)
			},
		},

		// scenario 5: Fail to create when an error was returned while fetching last movement
		{
			accountID: "123-123-123",
			newMov: &movement.Movement{
				ID:        movement.ID("1234567890"),
				AccountID: "123-123-123",
				IsDebit:   false,
				Amount:    1000,
			},
			lstMov:       nil,
			lstError:     errors.New("database error"),
			lockReturn:   nil,
			unlockReturn: nil,
			createError:  nil,
			asserts: func(t *testing.T, al *mock.AccountLocker, mf *mock.MovementFinder, mc *mock.MovementCreator, result error) {
				assert.EqualError(t, result, "database error")
				al.AssertExpectations(t)
				mf.AssertExpectations(t)
				mc.AssertNumberOfCalls(t, "Create", 0)
			},
		},

		// scenario 6: Fail to create due database error
		{
			accountID: "123-123-123",
			newMov: &movement.Movement{
				ID:        movement.ID("1234567890"),
				AccountID: "123-123-123",
				IsDebit:   false,
				Amount:    1000,
			},
			lstMov:       nil,
			lstError:     nil,
			lockReturn:   nil,
			unlockReturn: nil,
			createError:  errors.New("database error"),
			asserts: func(t *testing.T, al *mock.AccountLocker, mf *mock.MovementFinder, mc *mock.MovementCreator, result error) {
				assert.EqualError(t, result, "database error")
				al.AssertExpectations(t)
				mf.AssertExpectations(t)
				mc.AssertExpectations(t)
			},
		},
	}

	for _, test := range tests {
		al := &mock.AccountLocker{}
		al.On("Lock", test.accountID).Return(test.lockReturn)
		al.On("Unlock", test.accountID).Return(test.unlockReturn).Times(1)

		mf := &mock.MovementFinder{}
		mf.On("Last", test.accountID).Return(test.lstMov, test.lstError)

		mc := &mock.MovementCreator{}
		mc.On("Create", testifyMock.Anything).Return(test.createError)

		l := New(al, mf, mc)
		result := l.CreateMovement(context.Background(), *test.newMov)
		test.asserts(lts.T(), al, mf, mc, result)

		// Unlock should ALWAYS run!
		al.AssertNumberOfCalls(lts.T(), "Unlock", 1)
	}
}

func TestLedgerTestSuite(t *testing.T) {
	suite.Run(t, &LedgerTestSuite{})
}
