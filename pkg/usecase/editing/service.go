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

// Repository is the interface that wraps the methods
// that must be implemented by the repository in order
// for editing service to perform its job
//
//	- EditTag, EditExpense, EditActivity are the main editing methods
//
//	- FindTagByID is used to check if tags exist when editing expenses/activities
//
//	- FindTagByName is used to check for duplicate tags when editing tag
//
// 	- FindActivityByID is used to check if activity exists when editing expense
type Repository interface {
	EditTag(domain.Tag) error
	EditExpense(domain.Expense) error
	EditActivity(domain.Activity) error
	FindTagByID(domain.TagID) (domain.Tag, error)
	FindTagByName(string) (domain.Tag, error)
	FindActivityByID(domain.ActivityID) (domain.Activity, error)
}

// ErrTagNameDuplicate is returned when trying to edit a tag with a name that already exists in store
var ErrTagNameDuplicate error = errors.New("Tag name duplicate")
