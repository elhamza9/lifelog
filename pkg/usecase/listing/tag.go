package listing

import (
	"github.com/elhamza90/lifelog/pkg/domain"
)

// AllTags returns a list of all tags stored in the repo
func (srv Service) AllTags() (*[]domain.Tag, error) {
	return srv.repo.ListAllTags()
}
