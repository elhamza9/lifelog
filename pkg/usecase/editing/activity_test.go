package editing_test

import (
	"strings"
	"testing"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
	"github.com/elhamza90/lifelog/pkg/store"
)

func TestEditActivity(t *testing.T) {
	repo.Tags = map[domain.TagID]domain.Tag{
		1: {ID: 1, Name: "tag-1"},
	}
	repo.Activities = map[domain.ActivityID]domain.Activity{
		1: {
			ID:       1,
			Label:    "Test Activity",
			Time:     time.Now().AddDate(0, -1, 0),
			Duration: time.Duration(time.Hour),
		},
	}

	tests := map[string]struct {
		act         domain.Activity
		expectedErr error
	}{
		"Correct": {
			act: domain.Activity{
				ID:       1,
				Label:    "Edited Test Activity",
				Place:    "Somewhere",
				Tags:     []domain.Tag{{ID: 1}},
				Time:     time.Now().AddDate(0, 0, -1),
				Duration: time.Duration(time.Hour),
			},
			expectedErr: nil,
		},
		"Non Existing Activity": {
			act: domain.Activity{
				ID:       1232,
				Label:    "Edited Test Activity",
				Place:    "Somewhere",
				Tags:     []domain.Tag{{ID: 1}},
				Time:     time.Now().AddDate(0, 0, -1),
				Duration: time.Duration(time.Hour),
			},
			expectedErr: store.ErrActivityNotFound,
		},

		"Non-Existing Tag": {
			act: domain.Activity{
				ID:       1,
				Label:    "Edited Test Activity",
				Place:    "Somewhere",
				Tags:     []domain.Tag{{ID: 9898}},
				Time:     time.Now().AddDate(0, 0, -1),
				Duration: time.Duration(time.Hour),
			},
			expectedErr: store.ErrTagNotFound,
		},
		"Field Invalid (Time future)": {
			act: domain.Activity{
				ID:       1,
				Label:    "Edited Test Activity",
				Place:    "Somewhere",
				Tags:     []domain.Tag{{ID: 1}},
				Time:     time.Now().AddDate(0, 0, 1),
				Duration: time.Duration(time.Hour),
			},
			expectedErr: domain.ErrActivityTimeFuture,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := editor.EditActivity(test.act)
			if err != test.expectedErr {
				t.Fatalf("Expected Err: %v\nReturned Err: %v", test.expectedErr, err)
			}

			edited := repo.Activities[test.act.ID]
			expectedPlace := strings.ToLower(test.act.Place)
			if err == nil && edited.Place != expectedPlace {
				t.Fatalf("Expected: %v\nReturned: %v", expectedPlace, edited.Place)
			}
		})
	}
}
