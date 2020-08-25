package listing_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
)

func TestActivitiesByTime(t *testing.T) {
	now := time.Now()
	repo.Activities = &map[domain.ActivityID]domain.Activity{
		1: {
			ID:       1,
			Label:    "Activity last year",
			Time:     now.AddDate(-1, 0, 0),
			Duration: time.Duration(time.Hour),
		},
		2: {
			ID:       2,
			Label:    "Activity yesterday",
			Time:     now.AddDate(0, 0, -1),
			Duration: time.Duration(time.Hour),
		},
		3: {
			ID:       3,
			Label:    "Activity Last month",
			Time:     now.AddDate(0, -1, 0),
			Duration: time.Duration(time.Hour),
		},
		4: {
			ID:       4,
			Label:    "Activity 15 days ago",
			Time:     now.AddDate(0, 0, -15),
			Duration: time.Duration(time.Hour),
		},
	}
	tests := map[string]struct {
		minTime     time.Time
		expectedIDs []domain.ActivityID // In order !
		expectedErr error
	}{
		"Last 3 months": {minTime: now.AddDate(0, -3, 0), expectedIDs: []domain.ActivityID{2, 4, 3}, expectedErr: nil},
		"Future Date":   {minTime: now.AddDate(0, 0, 1), expectedIDs: []domain.ActivityID{}, expectedErr: domain.ErrActivityTimeFuture},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := lister.ActivitiesByTime(test.minTime)
			// Test Error
			if err != test.expectedErr {
				t.Fatalf("Expecting Error: %v\nReturned: %v", test.expectedErr, err)
			}
			// Test Result content and order
			errMsg := fmt.Sprintf("Expecting: %v\nReturned: %v", test.expectedIDs, res)
			if len(res) != len(test.expectedIDs) {
				t.Fatalf(errMsg)
			}
			for i, act := range res {
				if act.ID != test.expectedIDs[i] {
					t.Fatalf(errMsg)
				}
			}
		})
	}

}

func TestActivitiesByTag(t *testing.T) {
	now := time.Now()
	d := time.Duration(time.Minute * 45)
	repo.Tags = &map[domain.TagID]domain.Tag{
		1: {ID: 1, Name: "tag-1"},
		2: {ID: 2, Name: "tag-2"},
		3: {ID: 3, Name: "tag-3"},
	}

	repo.Activities = &map[domain.ActivityID]domain.Activity{
		1: {ID: 1, Label: "Act tag-1", Time: now.AddDate(0, 0, -3), Duration: d, Tags: []domain.Tag{{ID: 1}}},
		2: {ID: 2, Label: "Act tag-1/3", Time: now.AddDate(0, 0, -1), Duration: d, Tags: []domain.Tag{{ID: 3}, {ID: 1}}},
		3: {ID: 3, Label: "Act tag-2", Time: now.AddDate(0, 0, -1), Duration: d, Tags: []domain.Tag{{ID: 2}}},
		4: {ID: 4, Label: "Act tag-1/2", Time: now.AddDate(0, 0, -2), Duration: d, Tags: []domain.Tag{{ID: 2}, {ID: 1}}},
	}

	// Test Subcase: Non-Existing Tag
	t.Run("Non-Existing Tag", func(t *testing.T) {
		const nonExistingID domain.TagID = 988998
		_, err := lister.ActivitiesByTag(nonExistingID)
		expectedErr := domain.ErrTagNotFound
		failed := err != expectedErr
		if failed {
			t.Fatalf("Expected error: %v\nReturned Error: %v", expectedErr, err)
		}
	})

	// Test Subcase: Existing Tag
	t.Run("Existing Tag", func(t *testing.T) {
		const testID domain.TagID = 1
		res, err := lister.ActivitiesByTag(testID)
		if err != nil {
			t.Fatalf("Unexpected Error: %v", err)
		}
		expectedIDs := []domain.ActivityID{2, 4, 1}
		errMsg := fmt.Sprintf("Expected: %v\nReturned: %v", expectedIDs, res)
		if len(res) != len(expectedIDs) {
			t.Fatal(errMsg)
		}
		for i, act := range res {
			if act.ID != expectedIDs[i] {
				t.Fatal(errMsg)
			}
		}
	})

}
