package adding

import (
	"github.com/elhamza90/lifelog/internal/domain"
)

// NewActivity validates the activity and calls the repo to store it.
// It does the following checks:
//	- Check primitive fields are valid
//	- Check Tags exist in DB
func (srv Service) NewActivity(act domain.Activity) (domain.ActivityID, error) {
	// Check primitive fields are valid
	if err := act.Validate(); err != nil {
		return 0, err
	}

	// Check & Fetch Tags
	fetchedTags := []domain.Tag{}
	for _, t := range act.Tags {
		fetched, err := srv.repo.FindTagByID(t.ID)
		if err != nil {
			return 0, err
		}
		fetchedTags = append(fetchedTags, fetched)
	}
	act.Tags = fetchedTags

	return srv.repo.AddActivity(act)
}
