package adding

import "github.com/elhamza90/lifelog/pkg/domain"

// Service provides methods that create entities
// and call the given repository to store them
type Service struct {
	repo Repository
}

// NewService returns a new service with provided repository
func NewService(r Repository) Service {
	return Service{repo: r}
}

// Repository defines methods that must be implemented to store entities
type Repository interface {
	AddTag(domain.Tag) (domain.Tag, error)
	FindTagByName(string) (domain.Tag, error)
	AddExpense(domain.Expense) (domain.Expense, error)
	FindTagByID(domain.TagID) (domain.Tag, error)
	FindActivityByID(domain.ActivityID) (domain.Activity, error)
}
