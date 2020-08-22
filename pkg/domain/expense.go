package domain

import (
	"errors"
	"fmt"
)

// ExpenseID is a value-object representing Id of an expense
type ExpenseID uint

// Expense entity
type Expense struct {
	ID    ExpenseID
	Label string
	Value float32
	Unit  string
	Tags  []Tag
}

// Constants
const (
	ExpenseLabelMinLen uint = 3
	ExpenseLabelMaxLen uint = 50
	ExpenseUnitMinLen  uint = 3
	ExpenseUnitMaxLen  uint = 10
)

// Errors
var (
	ErrExpenseLabelLength = fmt.Errorf("Expense Label must be %d ~ %d characters long", ExpenseLabelMinLen, ExpenseLabelMaxLen)
	ErrExpenseValueErr    = errors.New("Expense Value must be strictly positive")
	ErrExpenseUnitLength  = fmt.Errorf("Expense Unit must %d ~ %d long", ExpenseUnitMinLen, ExpenseUnitMaxLen)
)
