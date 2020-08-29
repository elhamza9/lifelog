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

// Repository is the interface that wraps the methods that must be
// implemented by the repository in order for deleting service
// to perform its job
//
//	- DeleteExpense, DeleteExpensesByActivity, DeleteActivity, DeleteTag
//	  are the main methods to delete entities
//
//	- FindExpenseByID, FindActivityByID, FindTagByID are used to check
//	  for existance of entities before deleting them
//
//	- FindExpensesByActivity, FindExpensesByTag, FindActivitiesByTag are used
//	  to check if there are any things associated with tag before deleting it.
type Repository interface {
	DeleteTag(domain.TagID) error
	DeleteExpense(id domain.ExpenseID) error
	DeleteActivity(domain.ActivityID) error
	DeleteExpensesByActivity(domain.ActivityID) error
	FindActivityByID(domain.ActivityID) (domain.Activity, error)
	FindExpenseByID(id domain.ExpenseID) (domain.Expense, error)
	FindTagByID(domain.TagID) (domain.Tag, error)
	FindExpensesByActivity(domain.ActivityID) ([]domain.Expense, error)
	FindExpensesByTag(domain.TagID) ([]domain.Expense, error)
	FindActivitiesByTag(domain.TagID) ([]domain.Activity, error)
}
