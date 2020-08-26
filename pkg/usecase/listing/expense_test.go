package listing_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
	"github.com/elhamza90/lifelog/pkg/store"
)

func TestExpensesByTime(t *testing.T) {
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
			res, err := lister.ExpensesByTime(test.minTime)
			// Check returned error
			if err != test.expectedErr {
				t.Fatalf("Unexpected error: %v\nExpecting: %v", err, test.expectedErr)
			}
			errMsg := fmt.Sprintf("Expecting: %v\nGot: %v", test.expectedIDs, *res)
			// Check result length
			if len(*res) != len(test.expectedIDs) {
				t.Fatalf(errMsg)
			}
			// Check content and order of result
			for i, exp := range *res {
				if exp.ID != test.expectedIDs[i] {
					t.Fatalf(errMsg)
				}
			}
		})
	}
}

func TestExpensesByTag(t *testing.T) {
	now := time.Now()
	repo.Tags = &map[domain.TagID]domain.Tag{
		1: {ID: 1, Name: "tag-1"},
		2: {ID: 2, Name: "tag-2"},
		3: {ID: 3, Name: "tag-3"},
	}
	repo.Expenses = &map[domain.ExpenseID]domain.Expense{
		1: {
			ID:    1,
			Label: "6 months ago / tag-1 & tag-3",
			Time:  now.AddDate(0, -6, 0),
			Tags: []domain.Tag{
				{ID: 1, Name: "tag-1"},
				{ID: 3, Name: "tag-3"},
			},
		},
		2: {
			ID:    2,
			Label: "7 days ago / tag-2",
			Time:  now.AddDate(0, 0, -7),
			Tags: []domain.Tag{
				{ID: 2, Name: "tag-2"},
			},
		},
		3: {
			ID:    3,
			Label: "2 months ago / tag-3",
			Time:  now.AddDate(0, -2, 0),
			Tags: []domain.Tag{
				{ID: 3, Name: "tag-3"},
			},
		},
		4: {
			ID:    4,
			Label: "4 months ago / tag-2 & tag-3",
			Time:  now.AddDate(0, -4, 0),
			Tags: []domain.Tag{
				{ID: 2, Name: "tag-2"},
				{ID: 3, Name: "tag-3"},
			},
		},
	}

	// Test Subcase: Non-Existing Tag
	t.Run("Non-Existing Tag", func(t *testing.T) {
		const nonExistingID domain.TagID = 988998
		_, err := lister.ExpensesByTag(nonExistingID)
		expectedErr := store.ErrTagNotFound
		failed := err != expectedErr
		if failed {
			t.Fatalf("Expected error: %v\nReturned Error: %v", expectedErr, err)
		}
	})

	// Test Subcase: Existing Tag
	t.Run("Existing Tag", func(t *testing.T) {
		const testID domain.TagID = 3
		res, err := lister.ExpensesByTag(testID)
		if err != nil {
			t.Fatalf("Unexpected Error: %v", err)
		}
		expectedIDs := []domain.ExpenseID{3, 4, 1}
		errMsg := fmt.Sprintf("Expected: %v\nReturned: %v", expectedIDs, res)
		if len(*res) != len(expectedIDs) {
			t.Fatal(errMsg)
		}
		for i, exp := range *res {
			if exp.ID != expectedIDs[i] {
				t.Fatal(errMsg)
			}
		}
	})
}

func TestExpensesByActivity(t *testing.T) {

	act := domain.Activity{ID: 88}
	repo.Activities = &map[domain.ActivityID]domain.Activity{
		act.ID: act,
	}
	repo.Expenses = &map[domain.ExpenseID]domain.Expense{
		1: {
			ID:         1,
			Label:      "Test Exp",
			Value:      10,
			Unit:       "Dh",
			ActivityID: 0,
		},
		2: {
			ID:         2,
			Label:      "Test Exp",
			Value:      10,
			Unit:       "Dh",
			ActivityID: act.ID,
		},
		3: {
			ID:         3,
			Label:      "Test Exp",
			Value:      10,
			Unit:       "Dh",
			ActivityID: act.ID,
		},
		4: {
			ID:         4,
			Label:      "Test Exp",
			Value:      10,
			Unit:       "Dh",
			ActivityID: 0,
		},
	}
	// Test Subcase: Non existing Activity
	t.Run("Non-Existing Activity", func(t *testing.T) {
		nonExistantID := domain.ActivityID(989)
		_, err := lister.ExpensesByActivity(nonExistantID)
		expectdErr := store.ErrActivityNotFound
		if err != expectdErr {
			t.Fatalf("Expected Err: %v\nReturned Err: %v", store.ErrActivityNotFound, err)
		}
	})

	// Test Subcase existing Activity
	t.Run("Existing Activity", func(t *testing.T) {
		res, err := lister.ExpensesByActivity(act.ID)
		if err != nil {
			t.Fatalf("Unexpected Error: %v", err)
		}
		expectedIDs := []domain.ExpenseID{2, 3}
		errMsg := fmt.Sprintf("Expected: %v\nReturned: %v", expectedIDs, res)
		if len(*res) != len(expectedIDs) {
			t.Fatal(errMsg)
		}
		var found bool
		for _, exp := range *res {
			found = false
			for _, id := range expectedIDs {
				if exp.ID == id {
					found = true
					break
				}
			}
			if !found {
				t.Fatal(errMsg)
			}
		}
	})
}
