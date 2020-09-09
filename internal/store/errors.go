package store

import "errors"

// Errors
var (
	ErrTagNotFound      error = errors.New("Tag not found")
	ErrExpenseNotFound  error = errors.New("Expense Not Found")
	ErrActivityNotFound error = errors.New("Activity Not Found")
)
