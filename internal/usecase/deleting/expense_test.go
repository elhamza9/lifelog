package deleting_test

import (
	"testing"
	"time"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/elhamza90/lifelog/internal/store"
)

func TestDeleteExpense(t *testing.T) {
	repo.Expenses = map[domain.ExpenseID]domain.Expense{
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
		"Non-Existing Expense": {ID: 988998, expectedErr: store.ErrExpenseNotFound},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := deleter.Expense(test.ID)
			failed := err != test.expectedErr
			if failed {
				t.Fatalf("Expecting Error: %v\nReturned Error; %v", err, test.expectedErr)
			}
		})
	}
}

func TestDeleteActivityExpenses(t *testing.T) {
	now := time.Now()
	const testActID domain.ActivityID = 1
	repo.Activities = map[domain.ActivityID]domain.Activity{
		testActID: {ID: testActID, Label: "Test Act", Time: now.AddDate(0, 0, -1), Duration: time.Duration(time.Hour)},
	}

	// Subtests Definition
	tests := map[string]struct {
		actID       domain.ActivityID
		expectedErr error
	}{
		"Non-Existing Activity": {
			actID:       domain.ActivityID(98988),
			expectedErr: store.ErrActivityNotFound,
		},
		"Zero Activity ID": {
			actID:       domain.ActivityID(0),
			expectedErr: store.ErrActivityNotFound,
		},
		"Existing ID": {
			actID:       testActID,
			expectedErr: nil,
		},
	}
	// Subtests Execution
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := deleter.ActivityExpenses(test.actID)
			if err != test.expectedErr {
				t.Fatalf("\nExpected err: %v\nReturned err: %v", test.expectedErr, err)
			}
		})
	}
}
