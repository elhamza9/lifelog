package listing

import (
	"sort"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
)

// ActivitiesByTime returns activities with Time field greater than or equal to the given time.
// The returned activities are ordered from most recent to oldest
// It returns ErrActivityTimeFuture when given time is future
func (srv Service) ActivitiesByTime(t time.Time) ([]domain.Activity, error) {
	if t.After(time.Now()) {
		return []domain.Activity{}, domain.ErrActivityTimeFuture
	}
	res, err := srv.repo.FindActivitiesByTime(t)
	if err != nil {
		return []domain.Activity{}, err
	}
	// Sort
	sort.Slice(res, func(i, j int) bool {
		return (res)[i].Time.After((res)[j].Time)
	})
	return res, nil
}

// ActivitiesByTag returns expenses that have the tag with given ID
// in their Tags field
// The returned expenses are ordered from most recent to oldest
// It returns an error if tag with given ID is not found
func (srv Service) ActivitiesByTag(tid domain.TagID) ([]domain.Activity, error) {
	// Check if Tag exists
	if _, err := srv.repo.FindTagByID(tid); err != nil {
		return []domain.Activity{}, err
	}
	// Get Activities from Repo
	res, err := srv.repo.FindActivitiesByTag(tid)
	if err != nil {
		return []domain.Activity{}, err
	}
	// Sort using Time field descendent
	sort.Slice(res, func(i, j int) bool {
		elemI := (res)[i]
		elemJ := (res)[j]
		return elemI.Time.After(elemJ.Time)
	})
	return res, nil
}
