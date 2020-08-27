package adding

import (
	"github.com/elhamza90/lifelog/pkg/domain"
	"github.com/elhamza90/lifelog/pkg/usecase"
)

// NewTag creates the new tag and calls the service repository to store it.
//	- It transforms name to lowercase
//	- checks repo for tag with same name ( duplicate tags are not allowed )
func (srv Service) NewTag(name string) (domain.TagID, error) {
	// Create Tag
	t := domain.Tag{Name: name}

	// Check fields valid
	if err := t.Validate(); err != nil {
		return 0, err
	}

	// Check tag name is not duplicate
	if t, err := srv.repo.FindTagByName(t.Name); err != nil {
		return 0, err
	} else if len(t.Name) > 0 {
		return 0, usecase.ErrTagNameDuplicate
	}
	// Call repo to store it
	return srv.repo.AddTag(t)
}
