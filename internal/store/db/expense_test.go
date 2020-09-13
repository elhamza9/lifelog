package db_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/elhamza90/lifelog/internal/store"
	"github.com/elhamza90/lifelog/internal/store/db"
	"gorm.io/gorm"
)

func TestFindExpenseByID(t *testing.T) {
	defer clearDB()
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
		t.Fatalf("\nError while creating test expense:\n  %v", err)
	}
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
	defer clearDB()
	tags := []db.Tag{{ID: 1, Name: "test-tag-1"}, {ID: 2, Name: "test-tag-2"}}
	if err := grmDb.Create(&tags).Error; err != nil {
		t.Fatalf("\nError while creating Test Tags:\n  %v", err)
	}
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
	if err != nil {
		t.Fatalf("\nUnexpected Error: %v", err)
	}
	var created db.Expense
	if err := grmDb.Preload("Tags").First(&created, id).Error; err != nil {
		t.Fatalf("\nUnexpected Error while retrieving saved expense:\n  %v", err)
	}
	if len(created.Tags) != len(exp.Tags) {
		t.Fatalf("\nExpected %d Tags\nReturned %d Tags", len(exp.Tags), len(created.Tags))
	}
}

func TestFindExpensesByTime(t *testing.T) {
	// Create test 100 expenses:
	// one for each day starting from today going backward
	defer clearDB()
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
		t.Fatalf("\nError while creating test expenses:\n  %v", err)
	}
	// Test Get Expenses of last 5 days (Should be 6 expenses)
	minTime := now.AddDate(0, 0, -5)
	nbrExpectedExpenses := 6
	res, err := repo.FindExpensesByTime(minTime)
	if err != nil {
		t.Fatalf("\nUnexpected Error: %v", err)
	}
	if len(res) != nbrExpectedExpenses {
		t.Log(res)
		t.Fatalf("\nExpecting %d Expenses\nReturned %d expenses", nbrExpectedExpenses, len(res))
	}
	// Test Order by time
	for i, exp := range res {
		if i < len(res)-1 {
			if exp.Time.Before(res[i+1].Time) {
				t.Fatal("Expenses not ordered by time")
			}
		}
	}
}

func TestFindExpensesByTag(t *testing.T) {
	// Create test expenses & tags:
	defer clearDB()
	var (
		tag1 db.Tag = db.Tag{ID: 11, Name: "test-tag-1"}
		tag2 db.Tag = db.Tag{ID: 12, Name: "test-tag-2"}
		tag3 db.Tag = db.Tag{ID: 13, Name: "test-tag-3"}
	)
	if err := grmDb.Create(&[]db.Tag{tag1, tag2, tag3}).Error; err != nil {
		t.Fatalf("\nError while creating test tags:\n  %v", err)
	}
	now := time.Now()
	expenses := []db.Expense{
		{
			Label:      "Test Expense 1 ( Tag1, Tag3 )",
			Time:       now.AddDate(0, 0, -20),
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
		t.Fatalf("\nError while creating test expenses:\n  %v", err)
	}
	// Test Get Expenses of Tag 1
	res, err := repo.FindExpensesByTag(tag1.ID)
	if err != nil {
		t.Fatalf("\nUnexpected Error: %v", err)
	}
	expectedExpenses := [2]db.Expense{expenses[2], expenses[0]}
	if len(res) != len(expectedExpenses) {
		t.Fatalf("\nExpecting %d Expenses\nReturned %d expenses", len(expectedExpenses), len(res))
	}
	for i, exp := range res {
		if exp.ID != expectedExpenses[i].ID {
			t.Fatalf("\nExpecting expense ID %d in %d position, Got ID %d", expectedExpenses[i].ID, i+1, exp.ID)
		}
	}
}

func TestFindExpensesByActivity(t *testing.T) {
	// Create test expenses & activities:
	defer clearDB()
	var (
		act1 db.Activity = db.Activity{
			ID:       11,
			Label:    "Test Activity 1",
			Place:    "Somewhere",
			Time:     time.Now().AddDate(0, 0, -3),
			Duration: time.Duration(time.Hour),
		}
		act2 db.Activity = db.Activity{
			ID:       22,
			Label:    "Test Activity 2",
			Place:    "Somewhere",
			Time:     time.Now().AddDate(0, 0, -20),
			Duration: time.Duration(time.Hour),
		}
	)
	if err := grmDb.Create(&[]db.Activity{act1, act2}).Error; err != nil {
		t.Fatalf("\nError while creating test activities:\n  %v", err)
	}
	now := time.Now()
	expenses := []db.Expense{
		{
			Label:      "Test Expense 1 ( Act2 )",
			Time:       now.AddDate(0, 0, -20),
			Value:      10,
			Unit:       "eu",
			ActivityID: act2.ID,
			Tags:       []db.Tag{},
		},
		{
			Label:      "Test Expense 2 ( Act1)",
			Time:       now.AddDate(0, 0, -3),
			Value:      10,
			Unit:       "eu",
			ActivityID: act1.ID,
			Tags:       []db.Tag{},
		},
		{
			Label:      "Test Expense 3 ( Act2)",
			Time:       now.AddDate(0, 0, -15),
			Value:      10,
			Unit:       "eu",
			ActivityID: act2.ID,
			Tags:       []db.Tag{},
		},
	}
	if err := grmDb.Create(&expenses).Error; err != nil {
		t.Fatalf("\nError while creating test expenses:\n  %v", err)
	}
	// Test Get Expenses of Act 2
	res, err := repo.FindExpensesByActivity(act2.ID)
	if err != nil {
		t.Fatalf("\nUnexpected Error: %v", err)
	}
	expectedExpenses := [2]db.Expense{expenses[2], expenses[0]}
	if len(res) != len(expectedExpenses) {
		t.Fatalf("\nExpecting %d Expenses\nReturned %d expenses", len(expectedExpenses), len(res))
	}
	for i, exp := range res {
		if exp.ID != expectedExpenses[i].ID {
			t.Fatalf("\nExpecting expense ID %d in %d position, Got ID %d", expectedExpenses[i].ID, i+1, exp.ID)
		}
	}
}

func TestDeleteExpense(t *testing.T) {
	testFunc := func(id domain.ExpenseID, expectedErr error) string {
		if err := repo.DeleteExpense(id); err != expectedErr {
			return fmt.Sprintf("\nExpected Error: %v\nReturned Error: %v", expectedErr, err)
		}
		return ""
	}
	// Subcase: Non existing Expense
	t.Run("Non Existing Expense", func(t *testing.T) {
		if err := testFunc(domain.ExpenseID(24234234), store.ErrExpenseNotFound); err != "" {
			t.Fatal(err)
		}
	})
	// Subcase: Existing Expense
	t.Run("Existing Expense", func(t *testing.T) {
		// Create test expense
		defer clearDB()
		exp := db.Expense{
			ID:    123,
			Label: "Test Expense",
			Time:  time.Now().AddDate(0, 0, -1),
			Value: 10,
			Unit:  "eu",
		}
		if err := grmDb.Create(&exp).Error; err != nil {
			t.Fatalf("\nUnexpected Error while creating test expense:\n  %v", err)
		}
		// Test returned error
		if err := testFunc(exp.ID, nil); err != "" {
			t.Fatal(err)
		}
		// Test if expense in DB
		if err := grmDb.First(&db.Expense{}, exp.ID).Error; err != gorm.ErrRecordNotFound {
			t.Fatalf("\nExpected %v\nReturned: %v", gorm.ErrRecordNotFound, err)
		}
	})
}

func TestDeleteExpensesByActivity(t *testing.T) {
	// Create test expenses & activities
	defer clearDB()
	var (
		act1 db.Activity = db.Activity{
			ID:       11,
			Label:    "Test Activity 1",
			Place:    "Somewhere",
			Time:     time.Now().AddDate(0, 0, -3),
			Duration: time.Duration(time.Hour),
		}
		act2 db.Activity = db.Activity{
			ID:       22,
			Label:    "Test Activity 2",
			Place:    "Somewhere",
			Time:     time.Now().AddDate(0, 0, -20),
			Duration: time.Duration(time.Hour),
		}
	)
	if err := grmDb.Create(&[]db.Activity{act1, act2}).Error; err != nil {
		t.Fatalf("\nError while creating test activities:\n  %v", err)
	}
	now := time.Now()
	expenses := []db.Expense{
		{
			Label:      "Test Expense 1 ( Act2 )",
			Time:       now.AddDate(0, 0, -20),
			Value:      10,
			Unit:       "eu",
			ActivityID: act2.ID,
			Tags:       []db.Tag{},
		},
		{
			Label:      "Test Expense 2 ( Act1)",
			Time:       now.AddDate(0, 0, -3),
			Value:      10,
			Unit:       "eu",
			ActivityID: act1.ID,
			Tags:       []db.Tag{},
		},
		{
			Label:      "Test Expense 3 ( Act2)",
			Time:       now.AddDate(0, 0, -15),
			Value:      10,
			Unit:       "eu",
			ActivityID: act2.ID,
			Tags:       []db.Tag{},
		},
	}
	if err := grmDb.Create(&expenses).Error; err != nil {
		t.Fatalf("\nError while creating test expenses:\n  %v", err)
	}
	// Test Delete expense of activity Act2
	if err := repo.DeleteExpensesByActivity(act2.ID); err != nil {
		t.Fatalf("\nUnexpected Error: %v", err)
	}
	var res []db.Expense
	// Test if records still exist in DB
	if err := grmDb.Where("activity_id = ?", act2.ID).Find(&res).Error; err != nil {
		t.Fatalf("\nUnexpected Error while trying to retrieve records that should be deleted:\n  %v", err)
	}
	if len(res) != 0 {
		t.Log(res)
		t.Fatalf("\nRecords were not deleted")
	}
}

func TestEditExpense(t *testing.T) {
	testFunc := func(exp db.Expense, expectedErr error) string {
		if err := repo.EditExpense(exp.ToDomain()); err != expectedErr {
			return fmt.Sprintf("\nExpected Error: %v\nReturned Error: %v", expectedErr, err)
		}
		return ""
	}
	// Subcase: Non existing Expense
	t.Run("Non Existing Expense", func(t *testing.T) {
		exp := db.Expense{
			ID:         2343244, // non existing
			Label:      "Non Existing Edited",
			Time:       time.Now().AddDate(0, 0, -1),
			Value:      10,
			Unit:       "eu",
			ActivityID: 0,
			Tags:       []db.Tag{},
		}
		if err := testFunc(exp, store.ErrExpenseNotFound); err != "" {
			t.Fatal(err)
		}
	})
	/*
		// Subcase: Existing Expense
		t.Run("Existing Expense", func(t *testing.T) {
			// Create test expense
			defer clearDB()
			exp := db.Expense{
				ID:    123,
				Label: "Test Expense",
				Time:  time.Now().AddDate(0, 0, -1),
				Value: 10,
				Unit:  "eu",
			}
			if err := grmDb.Create(&exp).Error; err != nil {
				t.Fatalf("\nUnexpected Error while creating test expense:\n  %v", err)
			}
			// Test returned error
			if err := testFunc(exp.ID, nil); err != "" {
				t.Fatal(err)
			}
			// Test if expense in DB
			if err := grmDb.First(&db.Expense{}, exp.ID).Error; err != gorm.ErrRecordNotFound {
				t.Fatalf("\nExpected %v\nReturned: %v", gorm.ErrRecordNotFound, err)
			}
		})
	*/

}
