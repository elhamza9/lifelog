package adding

import (
	"errors"
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
func (srv Service) NewActivity(label string, place string, desc string, time time.Time, dur time.Duration, tags *[]domain.Tag) (domain.Activity, error) {
	return domain.Activity{}, errors.New("Not yet implemented")
}
