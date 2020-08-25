package deleting

import "github.com/elhamza90/lifelog/pkg/domain"

// DeleteExpense calls the repo to delete the expense with provided ID
// If expense with given ID does not exist returns error
func (srv Service) Expense(id domain.ExpenseID) error {
	// Check expense exist
	if _, err := srv.repo.FindExpenseByID(id); err != nil {
		return err
	}

	return srv.repo.DeleteExpense(id)
}
