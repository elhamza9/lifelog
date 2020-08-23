package adding_test

import (
	"strings"
	"testing"

	"github.com/elhamza90/lifelog/pkg/domain"
)

func TestNewExpense(t *testing.T) {
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
		val         float32
		unit        string
		tags        *[]domain.Tag
		expectedErr error
	}{
		"Correct": {
			label:       "my expense",
			val:         15.5,
			unit:        "Dh",
			tags:        &[]domain.Tag{{ID: 100001}, {ID: 100005}},
			expectedErr: nil,
		},
		"Short Label": {
			label:       "my",
			val:         15.5,
			unit:        "Dh",
			tags:        &[]domain.Tag{},
			expectedErr: domain.ErrExpenseLabelLength,
		},
		"Long Label": {
			label:       "my ver ver ver very very ver very very very very very very very very very long label",
			val:         15.5,
			unit:        "Dh",
			tags:        &[]domain.Tag{},
			expectedErr: domain.ErrExpenseLabelLength,
		},
		"Long Unit": {
			label:       "my expense",
			val:         15.5,
			unit:        "LongLongUnit",
			tags:        &[]domain.Tag{},
			expectedErr: domain.ErrExpenseUnitLength,
		},
		"Short Unit": {
			label:       "my expense",
			val:         15.5,
			unit:        "D",
			tags:        &[]domain.Tag{},
			expectedErr: domain.ErrExpenseUnitLength,
		},
		"Non-Existing Tag": {
			label:       "my expense",
			val:         15.5,
			unit:        "Dh",
			tags:        &[]domain.Tag{{ID: 200000}},
			expectedErr: domain.ErrTagNotFound,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resExp, err := service.NewExpense(test.label, test.val, test.unit, test.tags)
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
