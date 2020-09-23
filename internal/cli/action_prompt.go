package cli

import "github.com/manifoldco/promptui"

const (
	actionDelete   string = "Delete"
	actionEdit     string = "Edit"
	actionDetails  string = "Details"
	actionContinue string = "Continue"
	actionExit     string = "Exit"
)

// actionPrompt asks user what action does he want to perform on selected entity
func actionPrompt() (string, error) {
	prompt := promptui.Select{
		Label: "What do you want to do ?",
		Items: []string{actionDetails, actionEdit, actionDelete, actionContinue, actionExit},
	}
	_, op, err := prompt.Run()
	return op, err
}
