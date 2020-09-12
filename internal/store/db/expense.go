package db

import (
	"errors"
	"time"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/elhamza90/lifelog/internal/store"
	"gorm.io/gorm"
)

// FindExpenseByID returns expense with given ID
// It returns an error if expense not found
func (repo Repository) FindExpenseByID(id domain.ExpenseID) (domain.Expense, error) {
	var exp Expense
	err := repo.db.First(&exp, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = store.ErrExpenseNotFound
	}
	return exp.ToDomain(), err
}

// SaveExpense stores the given Expense in Db  and returns created expense's ID
func (repo Repository) SaveExpense(exp domain.Expense) (domain.ExpenseID, error) {
	tags := []Tag{}
	for _, t := range exp.Tags {
		tags = append(tags, Tag{ID: t.ID, Name: t.Name})
	}
	dbExp := Expense{
		ID:         exp.ID,
		Label:      exp.Label,
		Time:       exp.Time,
		Value:      exp.Value,
		Unit:       exp.Unit,
		ActivityID: exp.ActivityID,
		Tags:       tags,
	}
	res := repo.db.Create(&dbExp)
	return domain.ExpenseID(dbExp.ID), res.Error
}

// FindExpensesByTime returns expenses with Time field
// greater than or equal to provided time
func (repo Repository) FindExpensesByTime(t time.Time) ([]domain.Expense, error) {
	res := []Expense{}
	if err := repo.db.Where("time >= ?", t).Find(&res).Error; err != nil {
		return []domain.Expense{}, err
	}
	expenses := make([]domain.Expense, len(res))
	for i, exp := range res {
		expenses[i] = exp.ToDomain()
	}
	return expenses, nil
}
