package adding

import (
	"strings"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
)

// NewActivity creates a new activity and calls the repo to store it
// It does the following checks:
//	- Check Label length
//	- Check Place length
//	- Check Description length
//	- Check Time + Duration not future
//	- Check Tags exist in DB
func (srv Service) NewActivity(label string, place string, desc string, timeStart time.Time, dur time.Duration, tags *[]domain.Tag) (domain.Activity, error) {
	place = strings.ToLower(place)
	act := domain.Activity{
		Label:    label,
		Place:    place,
		Desc:     desc,
		Time:     timeStart,
		Duration: dur,
	}
	if err := act.Valid(); err != nil {
		return domain.Activity{}, err
	}

	// Check & Fetch Tags
	fetchedTags := []domain.Tag{}
	for _, t := range *tags {
		fetched, err := srv.repo.FindTagByID(t.ID)
		if err != nil {
			return domain.Activity{}, err
		}
		fetchedTags = append(fetchedTags, fetched)
	}
	act.Tags = fetchedTags

	return srv.repo.AddActivity(act)
}
