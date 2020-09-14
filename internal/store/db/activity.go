package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/elhamza90/lifelog/internal/store"
	"gorm.io/gorm"
)

// FindActivityByID returns activity with given ID.
// If none is found, returns error
func (repo Repository) FindActivityByID(id domain.ActivityID) (domain.Activity, error) {
	var act Activity
	err := repo.db.First(&act, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = store.ErrActivityNotFound
	}
	return act.ToDomain(), err
}

// SaveActivity stores the given activity in memory and returns created activity
func (repo Repository) SaveActivity(act domain.Activity) (domain.ActivityID, error) {
	tags := []Tag{}
	for _, t := range act.Tags {
		tags = append(tags, Tag{ID: t.ID, Name: t.Name})
	}
	dbAct := Activity{
		ID:    act.ID,
		Label: act.Label,
		Place: act.Place,
		Desc:  act.Desc,
		Time:  act.Time,
		Tags:  tags,
	}
	res := repo.db.Create(&dbAct)
	return domain.ActivityID(dbAct.ID), res.Error
}

// FindActivitiesByTime returns activities
// with Time field greater than or equal to the given time
func (repo Repository) FindActivitiesByTime(t time.Time) ([]domain.Activity, error) {
	res := []Activity{}
	if err := repo.db.Where("time >= ?", t).Order("time DESC").Find(&res).Error; err != nil {
		return []domain.Activity{}, err
	}
	activities := make([]domain.Activity, len(res))
	for i, exp := range res {
		activities[i] = exp.ToDomain()
	}
	return activities, nil
}

// FindActivitiesByTag returns activities that have the provided tag in their Tags field
func (repo Repository) FindActivitiesByTag(tid domain.TagID) ([]domain.Activity, error) {
	var tag Tag
	if err := repo.db.Preload("Activities", func(db *gorm.DB) *gorm.DB {
		return db.Order("activities.time DESC") // Order activities by time
	}).First(&tag, tid).Error; err != nil {
		return []domain.Activity{}, err
	}
	activities := make([]domain.Activity, len(tag.Activities))
	for i, exp := range tag.Activities {
		activities[i] = (*exp).ToDomain()
	}
	return activities, nil
}

// DeleteActivity removes activity with provided ID from memory
func (repo Repository) DeleteActivity(id domain.ActivityID) error {
	res := repo.db.Delete(&Activity{ID: id})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != 1 {
		return store.ErrActivityNotFound
	}
	return nil
}

// EditActivity edits given activity in memory
func (repo Repository) EditActivity(act domain.Activity) error {
	if err := repo.db.First(&Activity{}, act.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return store.ErrActivityNotFound
		}
		return err
	}
	tags := make([]Tag, len(act.Tags))
	for i, t := range act.Tags {
		tags[i] = Tag{ID: t.ID, Name: t.Name}
	}
	res := repo.db.Save(&Activity{
		ID:       act.ID,
		Label:    act.Label,
		Place:    act.Place,
		Desc:     act.Desc,
		Time:     act.Time,
		Duration: act.Duration,
		Tags:     tags,
	})
	if res.RowsAffected != 1 {
		return fmt.Errorf("%d Rows were affected", res.RowsAffected)
	}
	return res.Error
}
