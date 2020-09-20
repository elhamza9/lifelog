package cli

import (
	"errors"

	"github.com/manifoldco/promptui"
)

func loginPrompt() (string, error) {
	prompt := promptui.Prompt{
		Label: "Password",
		Validate: func(input string) error {
			if len(input) < 8 {
				return errors.New("Password too short")
			}
			return nil
		},
	}
	pass, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return pass, nil
}
