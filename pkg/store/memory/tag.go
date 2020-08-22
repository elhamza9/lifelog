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
