package adding

import (
	"strings"

	"github.com/elhamza90/lifelog/pkg/domain"
	"github.com/elhamza90/lifelog/pkg/usecase"
)

// NewTag creates the new tag and calls the service repository to store it.
//   - checks repo for tag with same name ( duplicate tags are not allowed )
func (srv Service) NewTag(name string) (domain.Tag, error) {
	t := domain.Tag{Name: name}
	if err := t.Valid(); err != nil {
		return domain.Tag{}, err
	}
	t.Name = strings.ToLower(t.Name)
	// Check tag name is not duplicate
	if t, err := srv.repo.FindTagByName(t.Name); err != nil {
		return domain.Tag{}, err
	} else if len(t.Name) > 0 {
		return domain.Tag{}, usecase.ErrTagNameDuplicate
	}
	// Create & Store
	return srv.repo.AddTag(t)
}
