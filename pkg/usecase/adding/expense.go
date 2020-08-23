package adding

import (
	"strings"

	"github.com/elhamza90/lifelog/pkg/domain"
)

// NewExpense creates the new expense and calls the service repository to store it.
// It does the following checks:
//	- Checks Label length
//	- Checks Value is strictly positive
//	- Checks Unit length
//	- Transform Unit to lowercase
//	- Checks Tags exist in Repo
func (srv Service) NewExpense(label string, value float32, unit string, tags *[]domain.Tag) (domain.Expense, error) {
	// Check Label
	if len(label) < domain.ExpenseLabelMinLen || len(label) > domain.ExpenseLabelMaxLen {
		return domain.Expense{}, domain.ErrExpenseLabelLength
	}
	// Check Unit and transform it to lowercase
	if len(unit) < domain.ExpenseUnitMinLen || len(unit) > domain.ExpenseUnitMaxLen {
		return domain.Expense{}, domain.ErrExpenseUnitLength
	}
	unit = strings.ToLower(unit)
	// Check & Fetch Tags
	fetchedTags := []domain.Tag{}
	for _, t := range *tags {
		fetched, err := srv.repo.FindTagByID(t.ID)
		if err != nil {
			return domain.Expense{}, err
		}
		fetchedTags = append(fetchedTags, fetched)
	}
	exp := domain.Expense{
		Label: label,
		Value: value,
		Unit:  unit,
		Tags:  fetchedTags,
	}
	return srv.repo.AddExpense(exp)
}
