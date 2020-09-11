package adding

import (
	"github.com/elhamza90/lifelog/internal/domain"
)

// Service provides methods that create entities
// and call the given repository to store them
type Service struct {
	repo Repository
}

// NewService returns a new adding service with provided repository
func NewService(r Repository) Service {
	return Service{repo: r}
}

// Repository is the interface that wraps the methods
// that must be implemented by the repository
// in order for adding service to perform its job.
//
// - SaveTag, SaveExpense and SaveActivity are the main
// methods to store the objects.
//
// - FindTagByName is used to check for duplicate tag names.
//
// - FindTagByID is used to check that tags exist when
//   creating an activity or an expense with tags.
//
// - FindActivityByID is used to check that an activity
//   exists when creating an expense.
type Repository interface {
	SaveTag(domain.Tag) (domain.TagID, error)
	SaveExpense(domain.Expense) (domain.ExpenseID, error)
	SaveActivity(domain.Activity) (domain.ActivityID, error)
	FindTagByName(string) (domain.Tag, error)
	FindTagByID(domain.TagID) (domain.Tag, error)
	FindActivityByID(domain.ActivityID) (domain.Activity, error)
}
