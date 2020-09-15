package editing

import (
	"errors"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/elhamza90/lifelog/internal/store"
)

// EditTag calls repo to edit the provided tag
func (srv Service) EditTag(t domain.Tag) error {
	// Check Tag valid
	if err := t.Validate(); err != nil {
		return err
	}
	// Check Tag exists
	if _, err := srv.repo.FindTagByID(t.ID); err != nil {
		return err
	}
	// Check tag name is not duplicate
	if t, err := srv.repo.FindTagByName(t.Name); (err != nil) && !errors.Is(err, store.ErrTagNotFound) {
		return err
	} else if len(t.Name) > 0 {
		return domain.ErrTagNameDuplicate
	}
	return srv.repo.EditTag(t)
}
