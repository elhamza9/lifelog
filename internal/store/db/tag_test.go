package db_test

import (
	"testing"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/elhamza90/lifelog/internal/store"
	"github.com/elhamza90/lifelog/internal/store/db"
)

func TestFindTagByID(t *testing.T) {
	// Create test Tag
	tag := db.Tag{ID: 546, Name: "test-tag"}
	grmDb.Create(&tag)
	defer grmDb.Delete(&tag)
	// Tests
	tests := map[string]struct {
		id          domain.TagID
		expectedErr error
	}{
		"Existing Tag":     {tag.ID, nil},
		"Non Existing Tag": {23423423, store.ErrTagNotFound},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := repo.FindTagByID(test.id); err != test.expectedErr {
				t.Fatalf("\nExpected Error: %v\nReturned Error: %v", test.expectedErr, err)
			}
		})
	}
}
