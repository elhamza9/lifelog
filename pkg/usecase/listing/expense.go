package listing

import (
	"sort"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
)

// FindExpensesByTime returns expenses with Time field greater than or equal to the given time.
// The returned expenses are ordered from most recent to oldest
// It returns ErrExpenseTimeFuture when given time is future
func (srv Service) FindExpensesByTime(t time.Time) ([]domain.Expense, error) {
	if t.After(time.Now()) {
		return []domain.Expense{}, domain.ErrExpenseTimeFuture
	}
	res, err := srv.repo.FindExpensesByTime(t)
	if err != nil {
		return []domain.Expense{}, err
	}
	// Sort using Time field descendent
	sort.Slice(*res, func(i, j int) bool {
		elemI := (*res)[i]
		elemJ := (*res)[j]
		return elemI.Time.After(elemJ.Time)
	})
	return *res, nil
}
