package movement

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MovementTestSuite struct {
	suite.Suite
}

func (ms *MovementTestSuite) TestMovementBalance() {
	tests := []struct {
		isDebit  bool
		amount   int
		prevBal  int
		expected int
	}{
		{
			isDebit:  false,
			amount:   10,
			prevBal:  0,
			expected: 10,
		},
		{
			isDebit:  false,
			amount:   10,
			prevBal:  10,
			expected: 20,
		},
		{
			isDebit:  false,
			amount:   10,
			prevBal:  20,
			expected: 30,
		},
		{
			isDebit:  true,
			amount:   10,
			prevBal:  10,
			expected: 0,
		},
		{
			isDebit:  true,
			amount:   10,
			prevBal:  20,
			expected: 10,
		},
		{
			isDebit:  true,
			amount:   100,
			prevBal:  20,
			expected: -80,
		},
	}

	for _, test := range tests {
		mov := &Movement{
			ID:               "",
			AccountID:        "",
			IsDebit:          test.isDebit,
			Amount:           test.amount,
			PreviousMovement: "",
			PreviousBalance:  test.prevBal,
			CreatedAt:        time.Time{},
		}
		assert.Equal(ms.T(), test.expected, mov.CurrentBalance())
	}
}

func TestMovementTestSuite(t *testing.T) {
	suite.Run(t, &MovementTestSuite{})
}
