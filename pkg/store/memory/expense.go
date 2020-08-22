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

// AddExpense stores the given Expense in memory  and returns created expense
func (repo Repository) AddExpense(exp domain.Expense) (domain.Expense, error) {
	exp.ID = generateRandomExpenseID()
	(*repo.Expenses)[exp.ID] = exp
	return exp, nil
}
