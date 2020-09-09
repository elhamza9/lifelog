package editing

import (
	"github.com/elhamza90/lifelog/internal/domain"
)

// EditExpense calls repo to update given expense
func (srv Service) EditExpense(exp domain.Expense) error {
	// Check primitive fields are valid
	if err := exp.Validate(); err != nil {
		return err
	}

	// Check Activity exists
	if exp.ActivityID > 0 {
		if _, err := srv.repo.FindActivityByID(exp.ActivityID); err != nil {
			return err
		}
	}

	// Check & Fetch Tags
	fetchedTags := []domain.Tag{}
	for _, t := range exp.Tags {
		fetched, err := srv.repo.FindTagByID(t.ID)
		if err != nil {
			return err
		}
		fetchedTags = append(fetchedTags, fetched)
	}
	exp.Tags = fetchedTags

	return srv.repo.EditExpense(exp)

}
