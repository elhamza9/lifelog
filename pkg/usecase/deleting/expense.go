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

// ActivityExpenses calls repo to delete all expenses belonging to
// provided activity
func (srv Service) ActivityExpenses(aid domain.ActivityID) error {
	// Check if activity with provided ID exists
	if _, err := srv.repo.FindActivityByID(aid); err != nil {
		return err
	}
	return srv.repo.DeleteExpensesByActivity(aid)
}
