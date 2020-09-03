package adding_test

import (
	"strings"
	"testing"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
	"github.com/elhamza90/lifelog/pkg/store"
)

func TestNewExpense(t *testing.T) {
	now := time.Now()
	// Init Repo with one activity to test checking if activity exist
	repo.Activities = map[domain.ActivityID]domain.Activity{
		100000: {ID: 100000, Label: "Test Activity", Time: now.AddDate(0, 0, -1), Duration: time.Duration(time.Hour)},
	}
	// Init Repo with some tags to test checking if tags exist
	repo.Tags = map[domain.TagID]domain.Tag{
		100000: {ID: 100000, Name: "tag-100000"},
		100001: {ID: 100001, Name: "tag-100001"},
		100002: {ID: 100002, Name: "tag-100002"},
		100003: {ID: 100003, Name: "tag-100003"},
		100004: {ID: 100004, Name: "tag-100004"},
		100005: {ID: 100005, Name: "tag-100005"},
	}

	// Sub-tests
	tests := map[string]struct {
		label       string
		time        time.Time
		val         float32
		unit        string
		tags        []domain.Tag
		activityID  domain.ActivityID
		expectedErr error
	}{
		"Correct-with-activity": {
			label:       "my expense",
			time:        now.AddDate(0, 0, -1),
			val:         15.5,
			unit:        "Dh",
			tags:        []domain.Tag{{ID: 100001}, {ID: 100005}},
			activityID:  100000,
			expectedErr: nil,
		},
		"Correct-without-activity": {
			label:       "my expense",
			time:        now.AddDate(0, 0, -1),
			val:         15.5,
			unit:        "Dh",
			tags:        []domain.Tag{{ID: 100001}, {ID: 100005}},
			activityID:  0,
			expectedErr: nil,
		},
		"Non-existing-Activity": {
			label:       "my expense",
			time:        now.AddDate(0, 0, -1),
			val:         15.5,
			unit:        "Dh",
			tags:        []domain.Tag{{ID: 100001}, {ID: 100005}},
			activityID:  98899889,
			expectedErr: store.ErrActivityNotFound,
		},
		"Zero value": {
			label:       "my expense",
			time:        now.AddDate(0, 0, -1),
			val:         0,
			unit:        "Dh",
			activityID:  100000,
			tags:        []domain.Tag{{ID: 100001}, {ID: 100005}},
			expectedErr: domain.ErrExpenseValue,
		},
		"Time Future": {
			label:       "my expense",
			time:        now.AddDate(0, 0, 1),
			val:         15.5,
			unit:        "Dh",
			activityID:  100000,
			tags:        []domain.Tag{{ID: 100001}, {ID: 100005}},
			expectedErr: domain.ErrExpenseTimeFuture,
		},
		"Short Label": {
			label:       "my",
			time:        now.AddDate(0, 0, -1),
			val:         15.5,
			unit:        "Dh",
			activityID:  100000,
			tags:        []domain.Tag{},
			expectedErr: domain.ErrExpenseLabelLength,
		},
		"Long Label": {
			label:       "my ver ver ver very very ver very very very very very very very very very long label",
			time:        now.AddDate(0, 0, -1),
			val:         15.5,
			unit:        "Dh",
			activityID:  100000,
			tags:        []domain.Tag{},
			expectedErr: domain.ErrExpenseLabelLength,
		},
		"Long Unit": {
			label:       "my expense",
			time:        now.AddDate(0, 0, -1),
			val:         15.5,
			unit:        "LongLongUnit",
			activityID:  100000,
			tags:        []domain.Tag{},
			expectedErr: domain.ErrExpenseUnitLength,
		},
		"Short Unit": {
			label:       "my expense",
			time:        now.AddDate(0, 0, -1),
			val:         15.5,
			unit:        "D",
			activityID:  100000,
			tags:        []domain.Tag{},
			expectedErr: domain.ErrExpenseUnitLength,
		},
		"Non-Existing Tag": {
			label:       "my expense",
			time:        now.AddDate(0, 0, -1),
			val:         15.5,
			unit:        "Dh",
			activityID:  100000,
			tags:        []domain.Tag{{ID: 200000}},
			expectedErr: store.ErrTagNotFound,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			exp := domain.Expense{
				Label:      test.label,
				Time:       test.time,
				Value:      test.val,
				Unit:       test.unit,
				ActivityID: test.activityID,
				Tags:       test.tags,
			}
			createdID, err := adder.NewExpense(exp)
			testFailed := err != test.expectedErr
			if testFailed {
				t.Fatalf("\nExpected Error: %v\nReturned Error: %v", test.expectedErr, err)
			}
			// If expense was added without errors:
			if err == nil {
				// Fetch created expense directly from repo
				createdExpense := repo.Expenses[createdID]
				// Check unit was transformed to lowercase
				expectedUnit := strings.ToLower(test.unit)
				if createdExpense.Unit != expectedUnit {
					t.Fatalf("\nExpected Unit: %s\nReturned Unit: %s", expectedUnit, createdExpense.Unit)
				}
				// Check Tags were populated
				if len(createdExpense.Tags) != len(test.tags) {
					t.Fatalf("\nCreated expense has number of tags different from number of tags given")
				}
			}
		})
	}
}
