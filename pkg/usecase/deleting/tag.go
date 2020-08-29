package deleting

import (
	"errors"

	"github.com/elhamza90/lifelog/pkg/domain"
)

var (
	//ErrTagHasExpenses is returned when tag to be deleted has expenses associated with it
	ErrTagHasExpenses error = errors.New("Tag can not be deleted because there are expenses associated with it")
	//ErrTagHasActivities is returned when tag to be deleted has activities associated with it
	ErrTagHasActivities error = errors.New("Tag can not be deleted because there are activities associated with it")
)

// Tag calls repo to remove Tag
// It does the following checks:
//	- Check if tag exists
//	- Check if there are any expenses/activities associated with tag
func (srv Service) Tag(id domain.TagID) error {
	// Check if Tag exists
	if _, err := srv.repo.FindTagByID(id); err != nil {
		return err
	}

	// Check if tag has expenses
	if res, err := srv.repo.FindExpensesByTag(id); err != nil {
		return err
	} else if len(res) > 0 {
		return ErrTagHasExpenses
	}

	// Check if tag has activities
	if res, err := srv.repo.FindActivitiesByTag(id); err != nil {
		return err
	} else if len(res) > 0 {
		return ErrTagHasActivities
	}

	return srv.repo.DeleteTag(id)
}
