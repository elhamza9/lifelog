package db

import (
	"errors"

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
	tags := []domain.Tag{}
	for _, t := range exp.Tags {
		tags = append(tags, domain.Tag{ID: t.ID, Name: t.Name})
	}
	return domain.Expense{
		ID:         exp.ID,
		Label:      exp.Label,
		Time:       exp.Time,
		Value:      exp.Value,
		Unit:       exp.Unit,
		ActivityID: exp.ActivityID,
		Tags:       tags,
	}, err
}
