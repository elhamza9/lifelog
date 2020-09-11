package db

import (
	"errors"
	"log"

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
	return t.ToDomain(), err
}

// FindTagByName searches for a tag with the given name and returns it.
// It returns an Empty Tag if not found.
func (repo Repository) FindTagByName(name string) (domain.Tag, error) {
	var t Tag
	err := repo.db.Where("name = ?", name).First(&t).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = store.ErrTagNotFound
	}
	return t.ToDomain(), err
}

// SaveTag stores the given Tag in db and returns created tag ID
func (repo Repository) SaveTag(t domain.Tag) (domain.TagID, error) {
	dbTag := Tag{Name: t.Name}
	res := repo.db.Create(&dbTag)
	return domain.TagID(dbTag.ID), res.Error
}

// FindAllTags returns all stored tags in db
func (repo Repository) FindAllTags() ([]domain.Tag, error) {
	var res []Tag
	if err := repo.db.Find(&res).Error; err != nil {
		return []domain.Tag{}, err
	}
	tags := []domain.Tag{}
	for _, t := range res {
		tags = append(tags, t.ToDomain())
	}
	return tags, nil
}

// DeleteTag deletes tag from db
func (repo Repository) DeleteTag(id domain.TagID) error {
	err := repo.db.Delete(&Tag{}, id).Error
	return err
}

// EditTag edits given tag in DB
func (repo Repository) EditTag(t domain.Tag) error {
	res := repo.db.Model(&Tag{ID: t.ID}).Updates(map[string]interface{}{"name": t.Name})
	log.Print(res.RowsAffected)
	return res.Error
}
