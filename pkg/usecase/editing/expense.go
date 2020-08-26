package editing

import (
	"strings"

	"github.com/elhamza90/lifelog/pkg/domain"
)

// Expense calls repo to update given expense
func (srv Service) Expense(exp domain.Expense) error {
	// Transform unit to lowecase
	exp.Unit = strings.ToLower(exp.Unit)
	if err := exp.Valid(); err != nil {
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
