package deleting_test

import (
	"testing"
	"time"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/elhamza90/lifelog/internal/store"
	"github.com/elhamza90/lifelog/internal/usecase/deleting"
)

func TestDeleteActivity(t *testing.T) {
	repo.Expenses = map[domain.ExpenseID]domain.Expense{
		1: {
			ID:         1,
			Label:      "Expense of activity 10",
			ActivityID: 10,
			Time:       time.Now().AddDate(0, 0, -1),
			Value:      10,
			Unit:       "Eu",
		},
	}
	repo.Activities = map[domain.ActivityID]domain.Activity{
		9: {
			ID:       9,
			Label:    "Exp",
			Time:     time.Now().AddDate(0, 0, -1),
			Duration: time.Duration(time.Hour),
		},
		10: {
			ID:       10,
			Label:    "Exp",
			Time:     time.Now().AddDate(0, 0, -1),
			Duration: time.Duration(time.Hour),
		},
	}

	tests := map[string]struct {
		ID          domain.ActivityID
		expectedErr error
	}{
		"Existing Activity":     {ID: 9, expectedErr: nil},
		"Non-Existing Activity": {ID: 988998, expectedErr: store.ErrActivityNotFound},
		"Activity with Expense": {ID: 10, expectedErr: deleting.ErrActivityHasExpenses},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := deleter.Activity(test.ID)
			failed := err != test.expectedErr
			if failed {
				t.Fatalf("\nExpected Error: %v\nReturned Error; %v", err, test.expectedErr)
			}
		})
	}
}
