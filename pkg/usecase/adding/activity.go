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
	now := time.Now()
	// Check Label Length
	if len(label) < domain.ActivityLabelMinLen || len(label) > domain.ActivityLabelMaxLen {
		return domain.Activity{}, domain.ErrActivityLabelLength
	}
	// Check Place Length
	if len(place) > domain.ActivityPlaceMaxLen {
		return domain.Activity{}, domain.ErrActivityPlaceLength
	}
	// Transform Place to Lowercase
	place = strings.ToLower(place)
	// Check Desc Length
	if len(desc) > domain.ActivityDescMaxLen {
		return domain.Activity{}, domain.ErrActivityDescLength
	}
	// Check TimeEnd not future
	if timeEnd := timeStart.Add(dur); timeEnd.After(now) {
		return domain.Activity{}, domain.ErrActivityTimeFuture
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

	act := domain.Activity{
		Label:    label,
		Place:    place,
		Desc:     desc,
		Time:     timeStart,
		Duration: dur,
		Tags:     fetchedTags,
	}
	return srv.repo.AddActivity(act)
}
