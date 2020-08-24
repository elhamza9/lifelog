package deleting_test

import (
	"testing"

	"github.com/elhamza90/lifelog/pkg/domain"
)

func TestDeleteExpense(t *testing.T) {
	repo.Expenses = &map[domain.ExpenseID]domain.Expense{
		9: {
			ID:         9,
			Label:      "Exp",
			Value:      10,
			Unit:       "Dh",
			ActivityID: 0,
			Tags:       []domain.Tag{},
		},
	}

	tests := map[string]struct {
		ID          domain.ExpenseID
		expectedErr error
	}{
		"Existing Expense":     {ID: 9, expectedErr: nil},
		"Non-Existing Expense": {ID: 988998, expectedErr: domain.ErrExpenseNotFound},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := service.DeleteExpense(test.ID)
			failed := err != test.expectedErr
			if failed {
				t.Fatalf("Expecting Error: %v\nReturned Error; %v", err, test.expectedErr)
			}
		})
	}
}
