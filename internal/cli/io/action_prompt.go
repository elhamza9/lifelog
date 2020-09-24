package io

import "github.com/manifoldco/promptui"

// Actions
const (
	ActionDelete   string = "Delete"
	ActionEdit     string = "Edit"
	ActionDetails  string = "Details"
	ActionContinue string = "Continue"
	ActionExit     string = "Exit"
)

// ActionPrompt asks user what action does he want to perform on selected entity
func ActionPrompt() (string, error) {
	prompt := promptui.Select{
		Label: "What do you want to do ?",
		Items: []string{ActionDetails, ActionEdit, ActionDelete, ActionContinue, ActionExit},
	}
	_, op, err := prompt.Run()
	return op, err
}
