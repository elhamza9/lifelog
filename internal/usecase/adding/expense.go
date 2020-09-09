package adding

import (
	"github.com/elhamza90/lifelog/internal/domain"
)

// NewExpense validates the new expense and calls the service repository to store it.
// It does the following checks:
//	- Check primitive fields are valid
//	- Check Activity with provided ActivityID exists
//	- Checks Tags exist and fetch them
func (srv Service) NewExpense(exp domain.Expense) (domain.ExpenseID, error) {

	// Check primitive fields are valid
	if err := exp.Validate(); err != nil {
		return 0, err
	}

	// Check Activity exists
	if exp.ActivityID > 0 {
		if _, err := srv.repo.FindActivityByID(exp.ActivityID); err != nil {
			return 0, err
		}
	}

	// Check & Fetch Tags
	fetchedTags := []domain.Tag{}
	for _, t := range exp.Tags {
		fetched, err := srv.repo.FindTagByID(t.ID)
		if err != nil {
			return 0, err
		}
		fetchedTags = append(fetchedTags, fetched)
	}
	exp.Tags = fetchedTags

	return srv.repo.AddExpense(exp)
}
