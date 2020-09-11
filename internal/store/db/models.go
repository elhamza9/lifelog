package db

import (
	"time"

	"github.com/elhamza90/lifelog/internal/domain"
)

// Tag Model
type Tag struct {
	ID        domain.TagID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName specifies the name of the table for the tag model
func (t Tag) TableName() string { return "tags" }

// ToDomain converts calling Tag to Domain Tag
func (t Tag) ToDomain() domain.Tag {
	return domain.Tag{
		ID:   t.ID,
		Name: t.Name,
	}
}

// Expense Model
type Expense struct {
	ID         domain.ExpenseID
	Label      string
	Time       time.Time
	Value      float32
	Unit       string
	ActivityID domain.ActivityID // Foreign Key
	Tags       []Tag             `gorm:"many2many:expense_tags;"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// ToDomain converts calling Expense to Domain Expense
func (exp Expense) ToDomain() domain.Expense {
	tags := []domain.Tag{}
	for _, t := range exp.Tags {
		tags = append(tags, t.ToDomain())
	}
	return domain.Expense{
		ID:         exp.ID,
		Label:      exp.Label,
		Time:       exp.Time,
		Value:      exp.Value,
		Unit:       exp.Unit,
		ActivityID: exp.ActivityID,
		Tags:       tags,
	}
}

// TableName specifies the name of the table for the expense model
func (exp Expense) TableName() string { return "expenses" }
