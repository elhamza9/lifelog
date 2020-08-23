package listing_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
)

func TestFindExpensesByTime(t *testing.T) {
	now := time.Now()
	repo.Expenses = &map[domain.ExpenseID]domain.Expense{
		1: {ID: 1, Label: "exp a day ago", Time: now.AddDate(0, 0, -1), Value: 10, Unit: "Eu"},
		2: {ID: 2, Label: "exp 1 month ago", Time: now.AddDate(0, -1, 0), Value: 10, Unit: "Eu"},
		3: {ID: 3, Label: "exp 1 year ago", Time: now.AddDate(-1, 0, 0), Value: 10, Unit: "Eu"},
		4: {ID: 4, Label: "exp 15 days ago", Time: now.AddDate(0, 0, -15), Value: 10, Unit: "Eu"},
		5: {ID: 5, Label: "exp 2 months ago", Time: now.AddDate(0, -2, 0), Value: 10, Unit: "Eu"},
	}

	tests := map[string]struct {
		minTime     time.Time
		expectedIDs []domain.ExpenseID // In descending order !
		expectedErr error
	}{
		"1 month ago": {
			minTime:     now.AddDate(0, -1, 0),
			expectedIDs: []domain.ExpenseID{1, 4, 2},
			expectedErr: nil,
		},
		"1 year ago": {
			minTime:     now.AddDate(-1, 0, 0),
			expectedIDs: []domain.ExpenseID{1, 4, 2, 5, 3},
			expectedErr: nil,
		},
		"Future": {
			minTime:     now.AddDate(1, 0, 0),
			expectedIDs: []domain.ExpenseID{},
			expectedErr: domain.ErrExpenseTimeFuture,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := service.FindExpensesByTime(test.minTime)
			// Check returned error
			if err != test.expectedErr {
				t.Fatalf("Unexpected error: %v\nExpecting: %v", err, test.expectedErr)
			}
			errMsg := fmt.Sprintf("Expecting: %v\nGot: %v", test.expectedIDs, res)
			// Check result length
			if len(res) != len(test.expectedIDs) {
				t.Fatalf(errMsg)
			}
			// Check content and order of result
			for i, exp := range res {
				if exp.ID != test.expectedIDs[i] {
					t.Fatalf(errMsg)
				}
			}
		})
	}
}
