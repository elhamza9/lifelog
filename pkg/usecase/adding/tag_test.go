package adding_test

import (
	"os"
	"strings"
	"testing"

	"github.com/elhamza90/lifelog/pkg/domain"
	"github.com/elhamza90/lifelog/pkg/store/memory"
	"github.com/elhamza90/lifelog/pkg/usecase/adding"
)

var service adding.Service
var repo memory.Repository

func TestMain(m *testing.M) {
	repo = memory.NewRepository()      // Work with In-Memory DB
	service = adding.NewService(&repo) // Passing by reference to change db when testing
	os.Exit(m.Run())
}

func TestNewTag(t *testing.T) {
	// Init Repo with a tag to test duplicate case
	repo.Tags = &map[domain.TagID]domain.Tag{
		100000: {Name: "duplicate-tag"},
	}

	// Sub-tests
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
		"Too Short":      {"my", domain.ErrTagNameTooShort},
		"Too Long":       {"myveryveryveryveryveryverylongtag", domain.ErrTagNameTooLong},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resTag, err := service.NewTag(test.name)
			testFailed := err != test.expectedErr
			if testFailed {
				var expectedErrStr string = "No Error"
				if test.expectedErr != nil {
					expectedErrStr = test.expectedErr.Error()
				}
				var errStr string = "No error"
				if err != nil {
					errStr = err.Error()
				}
				t.Fatalf("\nUnexpected: %s\nExpecting : %s", errStr, expectedErrStr)
			}
			// If no error was returned, check if stored tag's name was transformed to lowercase
			if err == nil && resTag.Name != strings.ToLower(test.name) {
				t.Fatalf("Tag Name was not transformed to Lowercase.")
			}
		})
	}
}
