package editing_test

import (
	"strings"
	"testing"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
	"github.com/elhamza90/lifelog/pkg/store"
)

func TestEditExpense(t *testing.T) {
	repo.Expenses = &map[domain.ExpenseID]domain.Expense{
		1: {
			ID:         1,
			Label:      "Test Expense",
			ActivityID: 0,
			Tags:       []domain.Tag{},
			Time:       time.Now().AddDate(0, -1, 0),
			Value:      10,
			Unit:       "Eu",
		},
	}
	repo.Tags = &map[domain.TagID]domain.Tag{
		1: {ID: 1, Name: "tag-1"},
	}
	repo.Activities = &map[domain.ActivityID]domain.Activity{
		1: {
			ID:       1,
			Label:    "Test Activity",
			Time:     time.Now().AddDate(0, -1, 0),
			Duration: time.Duration(time.Hour),
		},
	}

	tests := map[string]struct {
		exp         domain.Expense
		expectedErr error
	}{
		"Correct": {
			exp: domain.Expense{
				ID:         1,
				Label:      "Edited Test Expense",
				ActivityID: 1,
				Tags:       []domain.Tag{{ID: 1}},
				Time:       time.Now().AddDate(0, 0, -1),
				Value:      100,
				Unit:       "Dollar",
			},
			expectedErr: nil,
		},
		"Non-Existing ACtivity": {
			exp: domain.Expense{
				ID:         1,
				Label:      "Edited Test Expense",
				ActivityID: 9898,
				Tags:       []domain.Tag{{ID: 1}},
				Time:       time.Now().AddDate(0, 0, -1),
				Value:      100,
				Unit:       "Dollar",
			},
			expectedErr: store.ErrActivityNotFound,
		},
		"Non-Existing Tag": {
			exp: domain.Expense{
				ID:         1,
				Label:      "Edited Test Expense",
				ActivityID: 1,
				Tags:       []domain.Tag{{ID: 9898}},
				Time:       time.Now().AddDate(0, 0, -1),
				Value:      100,
				Unit:       "Dollar",
			},
			expectedErr: store.ErrTagNotFound,
		},
		"Field Invalid (Time future)": {
			exp: domain.Expense{
				ID:         1,
				Label:      "Edited Test Expense",
				ActivityID: 1,
				Tags:       []domain.Tag{{ID: 1}},
				Time:       time.Now().AddDate(0, 0, 1),
				Value:      100,
				Unit:       "Dollar",
			},
			expectedErr: domain.ErrExpenseTimeFuture,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := editter.Expense(test.exp)
			if err != test.expectedErr {
				t.Fatalf("Expected Err: %v\nReturned Err: %v", test.expectedErr, err)
			}
			edited := (*repo.Expenses)[test.exp.ID]
			if err == nil && (edited.Unit != strings.ToLower(test.exp.Unit) || edited.Label != test.exp.Label || edited.Value != test.exp.Value || edited.ActivityID != test.exp.ActivityID || len(edited.Tags) != len(test.exp.Tags)) {
				t.Fatalf("Expected: %v\nReturned: %v", test.exp, edited)
			}
		})
	}
}
