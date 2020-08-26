package editing

import (
	"errors"

	"github.com/elhamza90/lifelog/pkg/domain"
)

// Service provides methods that delete entities
type Service struct {
	repo Repository
}

// NewService returns a new service with provided repository
func NewService(r Repository) Service {
	return Service{repo: r}
}

// Repository defines methods that must be implemented to edit entities
type Repository interface {
	FindTagByID(domain.TagID) (domain.Tag, error)
	FindTagByName(string) (domain.Tag, error)
	EditTag(domain.Tag) (domain.Tag, error)
	EditExpense(domain.Expense) (domain.Expense, error)
	FindActivityByID(domain.ActivityID) (domain.Activity, error)
	EditActivity(domain.Activity) (domain.Activity, error)
}

// ErrTagNameDuplicate is returned when trying to edit a tag with a name that already exists in store
var ErrTagNameDuplicate error = errors.New("Tag name duplicate")
