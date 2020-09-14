package db_test

import (
	"testing"
	"time"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/elhamza90/lifelog/internal/store"
	"github.com/elhamza90/lifelog/internal/store/db"
)

func TestFindActivityByID(t *testing.T) {
	defer clearDB()
	// Create test Activity
	act := db.Activity{
		ID:       546,
		Label:    "test activity",
		Place:    "Somewhere",
		Desc:     "Details",
		Time:     time.Now(),
		Duration: time.Duration(time.Hour),
		Tags:     []db.Tag{},
	}
	if err := grmDb.Create(&act).Error; err != nil {
		t.Fatalf("\nError while creating test activity:\n  %v", err)
	}
	// Tests
	tests := map[string]struct {
		id          domain.ActivityID
		expectedErr error
	}{
		"Existing Activity":     {act.ID, nil},
		"Non Existing Activity": {23423423, store.ErrActivityNotFound},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := repo.FindActivityByID(test.id); err != test.expectedErr {
				t.Fatalf("\nExpected Error: %v\nReturned Error: %v", test.expectedErr, err)
			}
		})
	}
}

func TestSaveActivity(t *testing.T) {
	t.Fatal("Test not yet implemented")
}

func TestFindActivitiesByTime(t *testing.T) {
	t.Fatal("Test not yet implemented")
}

func TestFindActivitiesByTag(t *testing.T) {
	t.Fatal("Test not yet implemented")
}

func TestDeleteActivity(t *testing.T) {
	t.Fatal("Test not yet implemented")
}

func TestEditActivity(t *testing.T) {
	t.Fatal("Test not yet implemented")
}
