package db_test

import (
	"fmt"
	"math/rand"
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

func TestSaveExpense(t *testing.T) {
	// Create test Expense
	tags := []db.Tag{{ID: 1, Name: "test-tag-1"}, {ID: 2, Name: "test-tag-2"}}
	if err := grmDb.Create(&tags).Error; err != nil {
		t.Fatalf("Error while creating Test Tags:\n  %s", err)
	}
	defer grmDb.Where("1 = 1").Delete(&db.Tag{})
	exp := domain.Expense{
		ID:         546,
		Label:      "test expense",
		Time:       time.Now(),
		Value:      150,
		Unit:       "eu",
		ActivityID: 0,
		Tags:       []domain.Tag{{ID: tags[0].ID}},
	}
	// Test Save
	id, err := repo.SaveExpense(exp)
	defer grmDb.Where("1 = 1").Delete(&db.Expense{})
	if err != nil {
		t.Fatalf("Unexpected Error: %v", err)
	}
	var created db.Expense
	if err := grmDb.Preload("Tags").First(&created, id).Error; err != nil {
		t.Fatalf("Unexpectd Error: %v", err)
	}
	if len(created.Tags) != len(exp.Tags) {
		t.Fatalf("Expected %d Tags\nReturned %d Tags", len(exp.Tags), len(created.Tags))
	}
}

func TestFindExpensesByTime(t *testing.T) {
	// Create test 100 expenses:
	// one for each day starting from today going backward
	const nbrExpenses int = 100
	expenses := make([]db.Expense, nbrExpenses)
	now := time.Now()
	for i := 0; i < nbrExpenses; i++ {
		expenses[i] = db.Expense{
			Label:      fmt.Sprintf("Test Expense %d", i),
			Time:       now.AddDate(0, 0, -i),
			Value:      10,
			Unit:       "eu",
			ActivityID: 0,
			Tags:       []db.Tag{},
		}
	}
	// Shuffle expenses before saving them to DB to avoid getting them by insertion order
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(expenses), func(i, j int) { expenses[i], expenses[j] = expenses[j], expenses[i] })
	if err := grmDb.Create(&expenses).Error; err != nil {
		t.Fatalf("\nError while creating test expenses:\n  %s", err.Error())
	}
	defer grmDb.Where("1 = 1").Delete(&db.Expense{})
	// Test Get Expenses of last 15 days (Should be 15 expenses)
	minTime := now.AddDate(0, 0, -5)
	nbrExpectedExpenses := 6
	res, err := repo.FindExpensesByTime(minTime)
	if err != nil {
		t.Fatalf("Unexpected Error: %s", err.Error())
	}
	if len(res) != nbrExpectedExpenses {
		t.Logf("%v", res)
		t.Fatalf("\nExpecting %d Expenses\nReturned %d expenses", nbrExpectedExpenses, len(res))
	}
}
