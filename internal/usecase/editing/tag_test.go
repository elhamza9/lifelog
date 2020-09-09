package editing_test

import (
	"strings"
	"testing"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/elhamza90/lifelog/internal/store"
)

func TestEditTag(t *testing.T) {
	repo.Tags = map[domain.TagID]domain.Tag{
		1: {ID: 1, Name: "tag-1"},
		2: {ID: 2, Name: "duplicate"},
	}

	tests := map[string]struct {
		tag         domain.Tag
		expectedErr error
	}{
		"Non-Existing Tag": {
			tag:         domain.Tag{ID: 22, Name: "non-existing"},
			expectedErr: store.ErrTagNotFound,
		},
		"Duplicate Tag": {
			tag:         domain.Tag{ID: 1, Name: "Duplicate"},
			expectedErr: domain.ErrTagNameDuplicate,
		},
		"Existing Tag": {
			tag:         domain.Tag{ID: 1, Name: "Edited-tag-1"},
			expectedErr: nil,
		},
	}

	// Test Subcase: Non-existing Tag
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := editor.EditTag(test.tag)
			if err != test.expectedErr {
				t.Fatalf("Expected Err: %v\nReturned Err: %v", test.expectedErr, err)
			}
			expectedName := strings.ToLower(test.tag.Name)
			edited := repo.Tags[test.tag.ID]
			if err == nil && edited.Name != expectedName {
				t.Fatalf("Expected name: %s\nReturned name: %s", expectedName, edited.Name)
			}
		})
	}

}
