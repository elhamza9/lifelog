package adding

import (
	"errors"

	"github.com/elhamza90/lifelog/pkg/domain"
)

// NewExpense creates the new expense and calls the service repository to store it.
// It does the following checks:
//	- Checks Label length
//	- Checks Value is strictly positive
//	- Checks Unit length
//	- Checks Tags exist in Repo
func (srv Service) NewExpense(label string, value float32, unit string, tags *[]domain.Tag) (domain.Expense, error) {

	return domain.Expense{}, errors.New("Not yet implemented")
}
