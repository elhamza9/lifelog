package db

import (
	"errors"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/elhamza90/lifelog/internal/store"
	"gorm.io/gorm"
)

// FindTagByID searches for a tag with the given ID and returns it.
// It returns ErrTagNotFound if no tag was found.
func (repo Repository) FindTagByID(id domain.TagID) (domain.Tag, error) {
	var t Tag
	err := repo.db.First(&t, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = store.ErrTagNotFound
	}
	return domain.Tag{ID: t.ID, Name: t.Name}, err
}

// FindTagByName searches for a tag with the given name and returns it.
// It returns an Empty Tag if not found.
func (repo Repository) FindTagByName(name string) (domain.Tag, error) {
	var t Tag
	err := repo.db.Where("name = ?", name).First(&t).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = store.ErrTagNotFound
	}
	return domain.Tag{ID: t.ID, Name: t.Name}, err
}
