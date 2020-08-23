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
	ListAllTags() (*[]domain.Tag, error)
	FindExpensesByTime(time.Time) (*[]domain.Expense, error)
}
