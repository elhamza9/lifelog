package db_test

import (
	"fmt"
	"testing"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/elhamza90/lifelog/internal/store"
	"github.com/elhamza90/lifelog/internal/store/db"
	"gorm.io/gorm"
)

func TestFindTagByID(t *testing.T) {
	// Create test Tag
	defer clearDB()
	tag := db.Tag{ID: 546, Name: "test-tag"}
	if err := grmDb.Create(&tag).Error; err != nil {
		t.Fatalf("\nUnexpected Error while creating test tag:\n  %v", err)
	}
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
	defer clearDB()
	tag := db.Tag{ID: 546, Name: "test-tag"}
	if err := grmDb.Create(&tag).Error; err != nil {
		t.Fatalf("\nUnexpected Error while creating test tag:\n  %v", err)
	}
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

func TestSaveTag(t *testing.T) {
	defer clearDB()
	tag := domain.Tag{ID: 546, Name: "test-tag"}
	id, err := repo.SaveTag(tag)
	if err != nil {
		t.Fatalf("\nUnexpected Error: %v", err)
	}
	// Check tag was inserted in DB
	var created db.Tag
	if err := grmDb.First(&created, id).Error; err != nil {
		t.Fatalf("\nUnexpected Error: %v", err)
	}
	// Check tag fields values correspond to values provided
	if created.Name != tag.Name {
		t.Fatalf("Field Values of Tag dont correspond to provided values:\n\tCreated Tag: %v\n\tProvided Tag: %v", created, tag)
	}
}

func TestFindAllTags(t *testing.T) {
	// testFunc executes the repo method, checks error and length of result
	testFunc := func(nbrTags int) string {
		res, err := repo.FindAllTags()
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
		tags := []db.Tag{}
		for i := 0; i < nbrTags; i++ {
			tags = append(tags, db.Tag{
				Name: fmt.Sprintf("tag-%d", i),
			})
		}
		if err := grmDb.Create(&tags).Error; err != nil {
			t.Fatalf("Error creating test tags: %v", err)
		}
		// Test
		if err := testFunc(nbrTags); err != "" {
			t.Fatal(err)
		}
	})
}

func TestDeleteTag(t *testing.T) {
	// Create test Tag
	defer clearDB()
	tag := db.Tag{ID: 546, Name: "test-tag"}
	if err := grmDb.Create(&tag).Error; err != nil {
		t.Fatalf("\nUnexpected Error while creating test tag:\n  %v", err)
	}
	tests := map[string]struct {
		id          domain.TagID
		expectedErr error
	}{
		"Existing Tag":     {tag.ID, nil},
		"Non Existing Tag": {3432534, nil},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Test return value of method
			if err := repo.DeleteTag(test.id); err != test.expectedErr {
				t.Fatalf("\nExpected Error: %v\nReturned Error: %v", test.expectedErr, err)
			}
			// Test if record is in DB by trying to retrieve it.
			if err := grmDb.First(&db.Tag{}, test.id).Error; err != gorm.ErrRecordNotFound {
				t.Fatalf("\nExpected %v\nReturned: %v", gorm.ErrRecordNotFound, err)
			}
		})
	}
}

func TestEditTag(t *testing.T) {
	testFunc := func(tag domain.Tag) string {
		if err := repo.EditTag(tag); err != nil {
			return fmt.Sprintf("\nUnexpected Error: %v", err)
		}
		return ""
	}
	// Subcase Non Existing Tag
	t.Run("Non Existing Tag", func(t *testing.T) {
		nonExistingTag := domain.Tag{ID: 234234234, Name: "edited-non-existing"}
		if err := testFunc(nonExistingTag); err != "" {
			t.Fatal(err)
		}
		// Test if the tag was not saved when editing
		if err := grmDb.First(&db.Tag{}, nonExistingTag.ID).Error; err != gorm.ErrRecordNotFound {
			t.Fatalf("\nExpected %v\nReturned: %v", gorm.ErrRecordNotFound, err)
		}
	})
	// Subcase Existing Tag
	t.Run("Existing Tag", func(t *testing.T) {
		// Create test Tag
		defer clearDB()
		tag := db.Tag{ID: 546, Name: "test-tag"}
		if err := grmDb.Create(&tag).Error; err != nil {
			t.Fatalf("\nUnexpected Error while creating test tag:\n  %v", err)
		}
		editedTag := domain.Tag{ID: tag.ID, Name: "edited-tag"}
		if err := testFunc(editedTag); err != "" {
			t.Fatal(err)
		}
		// Test if the fields were updated
		var res db.Tag
		if err := grmDb.First(&res, tag.ID).Error; err != nil {
			t.Fatalf("\nUnexpected Error while retrieving updated tag: %v", err)
		}
		if res.Name != editedTag.Name {
			t.Fatalf("\nField values of edited activity were not fully updated:\n\t%v\n\t%v", res, editedTag)
		}
	})
}
