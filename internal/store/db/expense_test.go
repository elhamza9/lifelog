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

func TestFindExpensesByTag(t *testing.T) {
	// Create test expenses & tags:
	var (
		tag1 db.Tag = db.Tag{ID: 11, Name: "test-tag-1"}
		tag2 db.Tag = db.Tag{ID: 12, Name: "test-tag-2"}
		tag3 db.Tag = db.Tag{ID: 13, Name: "test-tag-3"}
	)
	if err := grmDb.Create(&[]db.Tag{tag1, tag2, tag3}).Error; err != nil {
		t.Fatalf("\nError while creating test tags:\n  %s", err.Error())
	}
	defer grmDb.Where("1 = 1").Delete(&db.Tag{})
	now := time.Now()
	expenses := []db.Expense{
		{
			Label:      "Test Expense 1 ( Tag1, Tag3 )",
			Time:       now.AddDate(0, 0, -10),
			Value:      10,
			Unit:       "eu",
			ActivityID: 0,
			Tags:       []db.Tag{tag1, tag3},
		},
		{
			Label:      "Test Expense 2 ( Tag2, Tag3 )",
			Time:       now.AddDate(0, 0, -3),
			Value:      10,
			Unit:       "eu",
			ActivityID: 0,
			Tags:       []db.Tag{tag2, tag3},
		},
		{
			Label:      "Test Expense 3 ( Tag1, Tag2 )",
			Time:       now.AddDate(0, 0, -15),
			Value:      10,
			Unit:       "eu",
			ActivityID: 0,
			Tags:       []db.Tag{tag1, tag2},
		},
	}
	if err := grmDb.Create(&expenses).Error; err != nil {
		t.Fatalf("\nError while creating test expenses:\n  %s", err.Error())
	}
	defer grmDb.Exec("DELETE FROM expense_tags")
	defer grmDb.Where("1 = 1").Delete(&db.Expense{})
	// Test Get Expenses of Tag 1
	nbrExpectedExpenses := 2
	res, err := repo.FindExpensesByTag(tag1.ID)
	if err != nil {
		t.Fatalf("Unexpected Error: %s", err.Error())
	}
	if len(res) != nbrExpectedExpenses {
		t.Fatalf("\nExpecting %d Expenses\nReturned %d expenses", nbrExpectedExpenses, len(res))
	}
}
