package memory

import (
	"math/rand"
	"time"

	"github.com/elhamza90/lifelog/pkg/domain"
	"github.com/elhamza90/lifelog/pkg/store"
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
	return domain.Expense{}, store.ErrExpenseNotFound
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

// FindExpensesByTag returns expenses that have the provided tag in their Tags field
func (repo Repository) FindExpensesByTag(tid domain.TagID) (*[]domain.Expense, error) {
	res := []domain.Expense{}
	for _, exp := range *repo.Expenses {
		for _, tag := range exp.Tags {
			if tag.ID == tid {
				res = append(res, exp)
			}
		}
	}
	return &res, nil
}

// FindExpensesByActivity returns expenses with ActivityID matching given id
func (repo Repository) FindExpensesByActivity(aid domain.ActivityID) (*[]domain.Expense, error) {
	res := []domain.Expense{}
	for _, exp := range *repo.Expenses {
		if exp.ActivityID == aid {
			res = append(res, exp)
		}
	}
	return &res, nil
}

// DeleteExpense deletes expense from memory
func (repo Repository) DeleteExpense(id domain.ExpenseID) error {
	if _, ok := (*repo.Expenses)[id]; !ok {
		return store.ErrExpenseNotFound
	}
	delete(*repo.Expenses, id)
	return nil
}

// DeleteExpensesByActivity deletes all expenses with given ActivityID
func (repo Repository) DeleteExpensesByActivity(aid domain.ActivityID) error {
	ids := []domain.ExpenseID{}
	for id, exp := range *repo.Expenses {
		if exp.ActivityID == aid {
			ids = append(ids, id)
		}
	}
	for _, id := range ids {
		delete(*repo.Expenses, id)
	}
	return nil
}
