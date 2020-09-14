package db

import (
	"time"

	"github.com/elhamza90/lifelog/internal/domain"
)

// FindActivityByID returns activity with given ID.
// If none is found, returns error
func (repo Repository) FindActivityByID(id domain.ActivityID) (domain.Activity, error) {
	return domain.Activity{}, errNotImplemented
}

// SaveActivity stores the given activity in memory and returns created activity
func (repo Repository) SaveActivity(act domain.Activity) (domain.ActivityID, error) {
	return 0, errNotImplemented
}

// FindActivitiesByTime returns activities
// with Time field greater than or equal to the given time
func (repo Repository) FindActivitiesByTime(t time.Time) ([]domain.Activity, error) {
	res := []domain.Activity{}
	return res, errNotImplemented
}

// FindActivitiesByTag returns actenses that have the provided tag in their Tags field
func (repo Repository) FindActivitiesByTag(tid domain.TagID) ([]domain.Activity, error) {
	res := []domain.Activity{}
	return res, errNotImplemented
}

// DeleteActivity removes activity with provided ID from memory
func (repo Repository) DeleteActivity(id domain.ActivityID) error {
	return errNotImplemented
}

// EditActivity edits given activity in memory
func (repo Repository) EditActivity(act domain.Activity) error {
	return errNotImplemented
}
