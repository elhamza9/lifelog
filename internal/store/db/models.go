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

// TableName specifies the name of the table for the expense model
func (exp Expense) TableName() string { return "expenses" }
