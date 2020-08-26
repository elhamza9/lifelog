package adding

import (
	"strings"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
)

// NewExpense creates the new expense and calls the service repository to store it.
// It does the following checks:
//	- Checks ActivityID exists (foreign key)
//	- Checks Label length
//	- Checks Time is not future
//	- Checks Value is strictly positive
//	- Checks Unit length
//	- Transform Unit to lowercase
//	- Checks Tags exist in Repo
func (srv Service) NewExpense(label string, t time.Time, value float32, unit string, activityID domain.ActivityID, tags *[]domain.Tag) (domain.Expense, error) {

	// Transform unit to lowecase
	unit = strings.ToLower(unit)
	// Create Expense
	exp := domain.Expense{
		Label: label,
		Time:  t,
		Value: value,
		Unit:  unit,
	}
	// Check validitity of fields
	if err := exp.Valid(); err != nil {
		return domain.Expense{}, err
	}

	// Check Activity exists
	if activityID > 0 {
		if _, err := srv.repo.FindActivityByID(activityID); err != nil {
			return domain.Expense{}, err
		}
	}
	exp.ActivityID = activityID

	// Check & Fetch Tags
	fetchedTags := []domain.Tag{}
	for _, t := range *tags {
		fetched, err := srv.repo.FindTagByID(t.ID)
		if err != nil {
			return domain.Expense{}, err
		}
		fetchedTags = append(fetchedTags, fetched)
	}
	exp.Tags = fetchedTags

	return srv.repo.AddExpense(exp)
}
