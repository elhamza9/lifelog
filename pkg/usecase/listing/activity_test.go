package listing_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
)

func TestFindActivitiesByTime(t *testing.T) {
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
			res, err := service.FindActivitiesByTime(test.minTime)
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
