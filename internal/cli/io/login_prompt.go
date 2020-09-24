package io

import (
	"errors"

	"github.com/manifoldco/promptui"
)

// LoginPrompt asks user for credentials
func LoginPrompt() (string, error) {
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
