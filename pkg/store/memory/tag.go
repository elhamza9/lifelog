package memory

import (
	"math/rand"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
)

func generateRandomTagID() domain.TagID {
	rand.Seed(time.Now().UnixNano())
	res := rand.Intn(10000)
	return domain.TagID(res)
}

// FindTagByID searches for a tag with the given ID and returns it.
// It returns ErrTagNotFound if no tag was found.
func (repo Repository) FindTagByID(id domain.TagID) (domain.Tag, error) {
	for _, t := range *(repo.Tags) {
		if t.ID == id {
			return t, nil
		}
	}
	return domain.Tag{}, domain.ErrTagNotFound
}

// FindTagByName searches for a tag with the given name and returns it.
// It returns an Empty Tag if not found.
func (repo Repository) FindTagByName(n string) (domain.Tag, error) {
	for _, t := range *(repo.Tags) {
		if t.Name == n {
			return t, nil
		}
	}
	return domain.Tag{}, nil
}

// AddTag stores the given Tag in memory  and returns created tag
func (repo Repository) AddTag(t domain.Tag) (domain.Tag, error) {
	t.ID = generateRandomTagID()
	(*repo.Tags)[t.ID] = t
	return t, nil
}

// ListAllTags returns all stored tags in memory
func (repo Repository) ListAllTags() (*[]domain.Tag, error) {
	tags := []domain.Tag{}
	for _, t := range *(repo.Tags) {
		tags = append(tags, t)
	}
	return &tags, nil
}
