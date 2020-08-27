package editing

import (
	"strings"

	"github.com/elhamza90/lifelog/pkg/domain"
)

// EditActivity calls repo to update given activity
func (srv Service) EditActivity(act domain.Activity) error {

	// Check Activity Exists
	if _, err := srv.repo.FindActivityByID(act.ID); err != nil {
		return err
	}

	// Transform unit to lowecase
	act.Place = strings.ToLower(act.Place)

	// Check primitive fields are valid
	if err := act.Valid(); err != nil {
		return err
	}

	// Check & Fetch Tags
	fetchedTags := []domain.Tag{}
	for _, t := range act.Tags {
		fetched, err := srv.repo.FindTagByID(t.ID)
		if err != nil {
			return err
		}
		fetchedTags = append(fetchedTags, fetched)
	}
	act.Tags = fetchedTags

	return srv.repo.EditActivity(act)
}
