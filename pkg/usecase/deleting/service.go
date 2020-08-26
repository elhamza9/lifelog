package deleting

import "github.com/elhamza90/lifelog/pkg/domain"

// Service provides methods that delete entities
type Service struct {
	repo Repository
}

// NewService returns a new service with provided repository
func NewService(r Repository) Service {
	return Service{repo: r}
}

// Repository defines methods that must be implemented to delete entities
type Repository interface {
	DeleteExpense(id domain.ExpenseID) error
	DeleteExpensesByActivity(domain.ActivityID) error
	DeleteActivity(domain.ActivityID) error
	DeleteTag(domain.TagID) error
	FindExpenseByID(id domain.ExpenseID) (domain.Expense, error)
	FindExpensesByActivity(domain.ActivityID) (*[]domain.Expense, error)
	FindActivityByID(domain.ActivityID) (domain.Activity, error)
	FindTagByID(domain.TagID) (domain.Tag, error)
	FindExpensesByTag(domain.TagID) (*[]domain.Expense, error)
	FindActivitiesByTag(domain.TagID) (*[]domain.Activity, error)
}
