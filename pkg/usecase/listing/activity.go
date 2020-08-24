package listing

import (
	"sort"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
)

// FindActivitiesByTime returns activities with Time field greater than or equal to the given time.
// The returned activities are ordered from most recent to oldest
// It returns ErrActivityTimeFuture when given time is future
func (srv Service) FindActivitiesByTime(t time.Time) ([]domain.Activity, error) {
	if t.After(time.Now()) {
		return []domain.Activity{}, domain.ErrActivityTimeFuture
	}
	res, err := srv.repo.FindActivitiesByTime(t)
	if err != nil {
		return []domain.Activity{}, err
	}
	// Sort
	sort.Slice(*res, func(i, j int) bool {
		return (*res)[i].Time.After((*res)[j].Time)
	})
	return *res, nil
}
