package memory

import (
	"math/rand"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
)

func generateRandomExpenseID() domain.ExpenseID {
	rand.Seed(time.Now().UnixNano())
	res := rand.Intn(10000)
	return domain.ExpenseID(res)
}

// FindExpenseByID returns expense with given ID
// It returns an error if expense not found
func (repo Repository) FindExpenseByID(id domain.ExpenseID) (domain.Expense, error) {
	for _, exp := range *repo.Expenses {
		if exp.ID == id {
			return exp, nil
		}
	}
	return domain.Expense{}, domain.ErrExpenseNotFound
}

// AddExpense stores the given Expense in memory  and returns created expense
func (repo Repository) AddExpense(exp domain.Expense) (domain.Expense, error) {
	exp.ID = generateRandomExpenseID()
	(*repo.Expenses)[exp.ID] = exp
	return exp, nil
}

// FindExpensesByTime returns expenses with Time
// field greater than or equal to provided tim
func (repo Repository) FindExpensesByTime(t time.Time) (*[]domain.Expense, error) {
	res := []domain.Expense{}
	for _, exp := range *repo.Expenses {
		if !exp.Time.Before(t) {
			res = append(res, exp)
		}
	}
	return &res, nil
}

// DeleteExpense deletes expense from memory
func (repo Repository) DeleteExpense(id domain.ExpenseID) error {
	if _, ok := (*repo.Expenses)[id]; !ok {
		return domain.ErrExpenseNotFound
	}
	delete(*repo.Expenses, id)
	return nil
}
