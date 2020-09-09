package deleting_test

import (
	"testing"
	"time"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/elhamza90/lifelog/internal/usecase/deleting"
)

func TestDeleteTag(t *testing.T) {
	repo.Tags = map[domain.TagID]domain.Tag{
		1: {ID: 1, Name: "tag-1"},
		2: {ID: 2, Name: "tag-2"},
		3: {ID: 3, Name: "tag-3"},
	}

	repo.Expenses = map[domain.ExpenseID]domain.Expense{
		1: {
			ID:    1,
			Label: "Expense Tag 2",
			Value: 10,
			Unit:  "Eu",
			Tags:  []domain.Tag{{ID: 2, Name: "tag-2"}},
		},
	}

	repo.Activities = map[domain.ActivityID]domain.Activity{
		1: {
			ID:       1,
			Label:    "Activity Tag 3",
			Time:     time.Now().AddDate(0, 0, -2),
			Duration: time.Duration(time.Hour),
			Tags: []domain.Tag{
				{ID: 3, Name: "tag-3"},
			},
		},
	}

	// Subtests definitions
	tests := map[string]struct {
		id          domain.TagID
		expectedErr error
	}{
		"Tag without activities/expenses": {
			id:          1,
			expectedErr: nil,
		},
		"Tag with expense": {
			id:          2,
			expectedErr: deleting.ErrTagHasExpenses,
		},
		"Tag with activity": {
			id:          3,
			expectedErr: deleting.ErrTagHasActivities,
		},
	}

	// Subtests Execution
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := deleter.Tag(test.id)
			if err != test.expectedErr {
				t.Fatalf("Expected Err: %v\nReturned Err: %v", test.expectedErr, err)
			}
		})
	}

}
