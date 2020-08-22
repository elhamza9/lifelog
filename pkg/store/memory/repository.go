package memory

// This package is used only for testing services

import "github.com/elhamza90/lifelog/pkg/domain"

// Repository manages data in memory using maps
type Repository struct {
	Tags *map[domain.TagID]domain.Tag
}

// NewRepository returns a new memory Repository with
// map pointers initialized to empty maps
func NewRepository() Repository {
	return Repository{
		Tags: &map[domain.TagID]domain.Tag{},
	}
}