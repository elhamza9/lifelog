package adding

import (
	"strings"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
)

// NewExpense creates the new expense and calls the service repository to store it.
// It does the following checks:
//	- Transform Unit to lowercase
//	- Check primitive fields are valid
//	- Check Activity with provided ActivityID exists
//	- Checks Tags exist and fetch them
func (srv Service) NewExpense(label string, t time.Time, value float32, unit string, activityID domain.ActivityID, tags *[]domain.Tag) (domain.ExpenseID, error) {

	// Transform unit to lowecase
	unit = strings.ToLower(unit)

	// Create Expense
	exp := domain.Expense{
		Label: label,
		Time:  t,
		Value: value,
		Unit:  unit,
	}

	// Check primitive fields are valid
	if err := exp.Valid(); err != nil {
		return 0, err
	}

	// Check Activity exists
	if activityID > 0 {
		if _, err := srv.repo.FindActivityByID(activityID); err != nil {
			return 0, err
		}
	}
	exp.ActivityID = activityID

	// Check & Fetch Tags
	fetchedTags := []domain.Tag{}
	for _, t := range *tags {
		fetched, err := srv.repo.FindTagByID(t.ID)
		if err != nil {
			return 0, err
		}
		fetchedTags = append(fetchedTags, fetched)
	}
	exp.Tags = fetchedTags

	return srv.repo.AddExpense(exp)
}
