package adding

import (
	"errors"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/elhamza90/lifelog/internal/store"
)

// NewTag validates tag and calls the service repository to store it.
//	- It transforms name to lowercase
//	- checks repo for tag with same name ( duplicate tags are not allowed )
func (srv Service) NewTag(t domain.Tag) (domain.TagID, error) {
	// Check fields valid
	if err := t.Validate(); err != nil {
		return 0, err
	}

	// Check tag name is not duplicate
	if t, err := srv.repo.FindTagByName(t.Name); (err != nil) && !errors.Is(err, store.ErrTagNotFound) {
		return 0, err
	} else if len(t.Name) > 0 {
		return 0, domain.ErrTagNameDuplicate
	}
	// Call repo to store it
	return srv.repo.SaveTag(t)
}
