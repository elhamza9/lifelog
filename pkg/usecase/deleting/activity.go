package deleting

import (
	"errors"

	"github.com/elhamza90/lifelog/pkg/domain"
)

// ErrActivityHasExpenses is returned when activity to be deleted has expenses associated with it
var ErrActivityHasExpenses error = errors.New("Activity can not be deleted because there are expenses associated with it")

// Activity deletes activity with provided ID.
// It does the following checks:
// 	- Check Activity Exists
//	- Check Activity has no expenses
func (srv Service) Activity(id domain.ActivityID) error {
	// Check Activity Exists
	if _, err := srv.repo.FindActivityByID(id); err != nil {
		return err
	}

	// Check activity has no expenses
	if res, err := srv.repo.FindExpensesByActivity(id); err != nil {
		return err
	} else if len(*res) > 0 {
		return ErrActivityHasExpenses
	}
	return srv.repo.DeleteActivity(id)
}
