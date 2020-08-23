package domain

import (
	"errors"
	"fmt"
	"time"
)

// ExpenseID is a value-object representing Id of an expense
type ExpenseID uint

// Expense entity
type Expense struct {
	ID    ExpenseID
	Label string
	Time  time.Time
	Value float32
	Unit  string
	Tags  []Tag
}

// Constants
const (
	ExpenseLabelMinLen int = 3
	ExpenseLabelMaxLen int = 50
	ExpenseUnitMinLen  int = 2
	ExpenseUnitMaxLen  int = 10
)

// Errors
var (
	ErrExpenseLabelLength = fmt.Errorf("Expense Label must be %d ~ %d characters long", ExpenseLabelMinLen, ExpenseLabelMaxLen)
	ErrExpenseValueErr    = errors.New("Expense Value must be strictly positive")
	ErrExpenseUnitLength  = fmt.Errorf("Expense Unit must %d ~ %d long", ExpenseUnitMinLen, ExpenseUnitMaxLen)
	ErrExpenseTimeFuture  = errors.New("Expense Time can not be future")
)
