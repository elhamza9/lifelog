package io

import "github.com/manifoldco/promptui"

// Actions
const (
	ActionDelete  string = "Delete"
	ActionEdit    string = "Edit"
	ActionDetails string = "Details"
	ActionExit    string = "Exit"
)

// CustomPrompt asks user to select an option out of given options
func CustomPrompt(question string, options []string) (string, error) {
	prompt := promptui.Select{
		Label: question,
		Items: options,
	}
	_, op, err := prompt.Run()
	return op, err
}

// ActionPrompt asks user what action does he want to perform on selected entity
func ActionPrompt() (string, error) {
	return CustomPrompt("What do you want to do ?", []string{ActionDetails, ActionEdit, ActionDelete, ActionExit})
}
