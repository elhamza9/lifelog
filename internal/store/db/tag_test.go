package db_test

import (
	"fmt"
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
	defer grmDb.Where("1 = 1").Delete(&db.Tag{})
	if err != nil {
		t.Fatalf("Unexpected Error: %v", err)
	}
	var created db.Tag
	if err := grmDb.First(&created, id).Error; err != nil {
		t.Fatalf("Unexpectd Error: %v", err)
	}
}

func TestListAllTags(t *testing.T) {
	// testFunc executes the repo method, checks error and length of result
	testFunc := func(nbrTags int) string {
		res, err := repo.ListAllTags()
		if err != nil {
			return fmt.Sprintf("\nUnexpected Error: %v", err)
		}
		if len(res) != nbrTags {
			return fmt.Sprintf("Expected 0 elements but got %d elements", len(res))
		}
		return ""
	}
	// Subcase Not Tags in DB
	t.Run("No Tags in DB", func(t *testing.T) {
		err := testFunc(0)
		if err != "" {
			t.Fatal(err)
		}
	})
	// Subcase multiple Tags in DB
	const nbrTags int = 100
	t.Run(fmt.Sprintf("%d Tags in DB", nbrTags), func(t *testing.T) {
		// Create some test tags
		defer grmDb.Where("1 = 1").Delete(&db.Tag{})
		for i := 0; i < nbrTags; i++ {
			grmDb.Create(&db.Tag{
				Name: fmt.Sprintf("tag-%d", i),
			})
		}
		if err := testFunc(nbrTags); err != "" {
			t.Fatal(err)
		}
	})
}
