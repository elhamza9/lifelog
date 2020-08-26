package editing

import (
	"strings"

	"github.com/elhamza90/lifelog/pkg/domain"
	"github.com/elhamza90/lifelog/pkg/usecase"
)

// Tag calls repo to edit the provided tag
func (srv Service) Tag(t domain.Tag) (domain.Tag, error) {
	// Transform name to lowercase
	t.Name = strings.ToLower(t.Name)
	// Check Tag valid
	if err := t.Valid(); err != nil {
		return domain.Tag{}, err
	}
	// Check Tag exists
	if _, err := srv.repo.FindTagByID(t.ID); err != nil {
		return domain.Tag{}, err
	}
	// Check tag name is not duplicate
	if t, err := srv.repo.FindTagByName(t.Name); err != nil {
		return domain.Tag{}, err
	} else if len(t.Name) > 0 {
		return domain.Tag{}, usecase.ErrTagNameDuplicate
	}
	return srv.repo.EditTag(t)
}
