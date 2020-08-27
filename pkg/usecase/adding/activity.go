package adding

import (
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
)

// NewActivity creates a new activity and calls the repo to store it.
// It does the following checks:
//	- Check primitive fields are valid
//	- Check Tags exist in DB
func (srv Service) NewActivity(label string, place string, desc string, timeStart time.Time, dur time.Duration, tags *[]domain.Tag) (domain.ActivityID, error) {
	// Create Activity
	act := domain.Activity{
		Label:    label,
		Place:    place,
		Desc:     desc,
		Time:     timeStart,
		Duration: dur,
	}

	// Check primitive fields are valid
	if err := act.Validate(); err != nil {
		return 0, err
	}

	// Check & Fetch Tags
	fetchedTags := []domain.Tag{}
	for _, t := range *tags {
		fetched, err := srv.repo.FindTagByID(t.ID)
		if err != nil {
			return 0, err
		}
		fetchedTags = append(fetchedTags, fetched)
	}
	act.Tags = fetchedTags

	return srv.repo.AddActivity(act)
}
