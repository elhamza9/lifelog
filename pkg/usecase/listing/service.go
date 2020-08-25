package listing

import (
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
)

// Service provides methods that list entities
type Service struct {
	repo Repository
}

// NewService returns a new listing service with provided repository
func NewService(r Repository) Service {
	return Service{repo: r}
}

// Repository defines methods that must be implemented to list entities
type Repository interface {
	FindTagByID(domain.TagID) (domain.Tag, error)
	ListAllTags() (*[]domain.Tag, error)
	FindExpensesByTime(time.Time) (*[]domain.Expense, error)
	FindExpensesByTag(domain.TagID) (*[]domain.Expense, error)
	FindActivitiesByTag(domain.TagID) (*[]domain.Activity, error)
	FindActivitiesByTime(time.Time) (*[]domain.Activity, error)
}
