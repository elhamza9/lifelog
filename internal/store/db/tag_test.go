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

func TestFindTagByName(t *testing.T) {
	// Create test Tag
	tag := db.Tag{ID: 546, Name: "test-tag"}
	grmDb.Create(&tag)
	defer grmDb.Delete(&tag)
	// Tests
	tests := map[string]struct {
		name        string
		expectedErr error
	}{
		"Existing Tag":     {tag.Name, nil},
		"Non Existing Tag": {"some-non-existing-tag", store.ErrTagNotFound},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := repo.FindTagByName(test.name); err != test.expectedErr {
				t.Fatalf("\nExpected Error: %v\nReturned Error: %v", test.expectedErr, err)
			}
		})
	}
}

func TestAddTag(t *testing.T) {
	// Create test Tag
	tag := domain.Tag{ID: 546, Name: "test-tag"}
	id, err := repo.AddTag(tag)
	if err != nil {
		t.Fatalf("Unexpected Error: %v", err)
	}
	var created db.Tag
	if err := grmDb.First(&created, id).Error; err != nil {
		t.Fatalf("Unexpectd Error: %v", err)
	}
}
