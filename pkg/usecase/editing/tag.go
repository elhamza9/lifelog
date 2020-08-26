package editing

import (
	"strings"

	"github.com/elhamza90/lifelog/pkg/domain"
)

// Tag calls repo to edit the provided tag
func (srv Service) Tag(t domain.Tag) (domain.Tag, error) {
	// Check Tag exists
	if _, err := srv.repo.FindTagByID(t.ID); err != nil {
		return domain.Tag{}, err
	}
	t.Name = strings.ToLower(t.Name)

	return srv.repo.EditTag(t)
}
