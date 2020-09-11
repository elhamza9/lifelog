package db_test

import (
	"testing"
	"time"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/elhamza90/lifelog/internal/store"
	"github.com/elhamza90/lifelog/internal/store/db"
)

func TestFindExpenseByID(t *testing.T) {
	// Create test Expense
	exp := db.Expense{
		ID:         546,
		Label:      "test expense",
		Time:       time.Now(),
		Value:      150,
		Unit:       "eu",
		ActivityID: 0,
		Tags:       []db.Tag{},
	}
	if err := grmDb.Create(&exp).Error; err != nil {
		t.Fatalf("\nError while creating test expense:\n  %s", err)
	}
	defer grmDb.Delete(&exp)
	// Tests
	tests := map[string]struct {
		id          domain.ExpenseID
		expectedErr error
	}{
		"Existing Expense":     {exp.ID, nil},
		"Non Existing Expense": {23423423, store.ErrExpenseNotFound},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := repo.FindExpenseByID(test.id); err != test.expectedErr {
				t.Fatalf("\nExpected Error: %v\nReturned Error: %v", test.expectedErr, err)
			}
		})
	}
}
