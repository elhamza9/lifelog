package adding_test

import (
	"strings"
	"testing"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
)

func TestNewExpense(t *testing.T) {
	// Init Repo with some activities to test checking if activity exist
	repo.Activities = &map[domain.ActivityID]domain.Activity{
		100000: {ID: 100000, Label: "Test Activity", Time: time.Now().AddDate(0, 0, -1), Duration: time.Duration(time.Hour)},
	}
	// Init Repo with some tags to test checking if tags exist
	repo.Tags = &map[domain.TagID]domain.Tag{
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
		tags        *[]domain.Tag
		activityID  domain.ActivityID
		expectedErr error
	}{
		"Correct-with-activity": {
			label:       "my expense",
			time:        time.Now().AddDate(0, 0, -1),
			val:         15.5,
			unit:        "Dh",
			tags:        &[]domain.Tag{{ID: 100001}, {ID: 100005}},
			activityID:  100000,
			expectedErr: nil,
		},
		"Correct-without-activity": {
			label:       "my expense",
			time:        time.Now().AddDate(0, 0, -1),
			val:         15.5,
			unit:        "Dh",
			tags:        &[]domain.Tag{{ID: 100001}, {ID: 100005}},
			activityID:  0,
			expectedErr: nil,
		},
		"Non-existing-Activity": {
			label:       "my expense",
			time:        time.Now().AddDate(0, 0, -1),
			val:         15.5,
			unit:        "Dh",
			tags:        &[]domain.Tag{{ID: 100001}, {ID: 100005}},
			activityID:  98899889,
			expectedErr: domain.ErrActivityNotFound,
		},

		"Zero value": {
			label:       "my expense",
			time:        time.Now().AddDate(0, 0, -1),
			val:         0,
			unit:        "Dh",
			activityID:  100000,
			tags:        &[]domain.Tag{{ID: 100001}, {ID: 100005}},
			expectedErr: domain.ErrExpenseValue,
		},

		"Time Future": {
			label:       "my expense",
			time:        time.Now().AddDate(0, 0, 1),
			val:         15.5,
			unit:        "Dh",
			activityID:  100000,
			tags:        &[]domain.Tag{{ID: 100001}, {ID: 100005}},
			expectedErr: domain.ErrExpenseTimeFuture,
		},

		"Short Label": {
			label:       "my",
			time:        time.Now().AddDate(0, 0, -1),
			val:         15.5,
			unit:        "Dh",
			activityID:  100000,
			tags:        &[]domain.Tag{},
			expectedErr: domain.ErrExpenseLabelLength,
		},
		"Long Label": {
			label:       "my ver ver ver very very ver very very very very very very very very very long label",
			time:        time.Now().AddDate(0, 0, -1),
			val:         15.5,
			unit:        "Dh",
			activityID:  100000,
			tags:        &[]domain.Tag{},
			expectedErr: domain.ErrExpenseLabelLength,
		},
		"Long Unit": {
			label:       "my expense",
			time:        time.Now().AddDate(0, 0, -1),
			val:         15.5,
			unit:        "LongLongUnit",
			activityID:  100000,
			tags:        &[]domain.Tag{},
			expectedErr: domain.ErrExpenseUnitLength,
		},
		"Short Unit": {
			label:       "my expense",
			time:        time.Now().AddDate(0, 0, -1),
			val:         15.5,
			unit:        "D",
			activityID:  100000,
			tags:        &[]domain.Tag{},
			expectedErr: domain.ErrExpenseUnitLength,
		},
		"Non-Existing Tag": {
			label:       "my expense",
			time:        time.Now().AddDate(0, 0, -1),
			val:         15.5,
			unit:        "Dh",
			activityID:  100000,
			tags:        &[]domain.Tag{{ID: 200000}},
			expectedErr: domain.ErrTagNotFound,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resExp, err := service.NewExpense(test.label, test.time, test.val, test.unit, test.activityID, test.tags)
			testFailed := err != test.expectedErr
			var expectedErrStr string = "No Error"
			if test.expectedErr != nil {
				expectedErrStr = test.expectedErr.Error()
			}
			var errStr string = "No Error"
			if err != nil {
				errStr = err.Error()
			}
			if testFailed {
				t.Fatalf("\nExpecting: %s\nBut Got: %s", expectedErrStr, errStr)
			}
			if err == nil {
				if resExp.Unit != strings.ToLower(test.unit) {
					t.Fatalf("Expense Unit should be transformed to Lower case.\nGot: %s", resExp.Unit)
				}
				if len(resExp.Tags) != len((*test.tags)) {
					t.Fatalf("Created expense has number of tags different from number of tags given")
				}
			}

		})
	}
}
