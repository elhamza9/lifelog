package adding_test

import (
	"strings"
	"testing"

	"github.com/elhamza90/lifelog/pkg/domain"
)

func TestNewTag(t *testing.T) {
	// Init Repo with a tag to test duplicate case
	repo.Tags = map[domain.TagID]domain.Tag{
		100000: {Name: "duplicate-tag"},
	}

	// Sub-tests Definitions
	tests := map[string]struct {
		name        string
		expectedErr error
	}{
		"Correct":        {"my-TAG_1", nil},
		"Duplicate":      {"duplicate-tag", domain.ErrTagNameDuplicate},
		"Spaces":         {"my tag", domain.ErrTagNameInvalidCharacters},
		"Special Char &": {"my-tag&", domain.ErrTagNameInvalidCharacters},
		"Special Char %": {"my-tag%", domain.ErrTagNameInvalidCharacters},
		"Special Char *": {"my-tag*", domain.ErrTagNameInvalidCharacters},
		"Too Short":      {"my", domain.ErrTagNameLen},
		"Too Long":       {"myveryveryveryveryveryverylongtag", domain.ErrTagNameLen},
	}

	// Sub-tests execution
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tag := domain.Tag{Name: test.name}
			createdID, err := adder.NewTag(tag)
			testFailed := err != test.expectedErr
			if testFailed {
				t.Fatalf("\nExpected Error: %v\nReturned Error: %v", test.expectedErr, err)
			}
			// If no error was returned, check if stored tag's name was transformed to lowercase
			if err == nil {
				// Fetch created tag from repo
				createdTag := repo.Tags[createdID]
				expectedName := strings.ToLower(test.name)
				if createdTag.Name != expectedName {
					t.Fatalf("\nExpected Tag Name: %s\nReturned Tag Name: %s", expectedName, createdTag.Name)
				}
			}
		})
	}
}
