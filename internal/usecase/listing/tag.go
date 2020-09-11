package listing

import (
	"github.com/elhamza90/lifelog/internal/domain"
)

// AllTags returns a list of all tags stored in the repo
func (srv Service) AllTags() ([]domain.Tag, error) {
	return srv.repo.FindAllTags()
}

// GetTagByID returns a tag ith given ID
func (srv Service) GetTagByID(id domain.TagID) (domain.Tag, error) {
	return srv.repo.FindTagByID(id)
}
