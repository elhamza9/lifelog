package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

// yesNoPrompt asks user a Yes or No question and returns boolean result
func yesNoPrompt(question string, defaultAction string) (res bool, err error) {
	defaultActionSmallCase := strings.ToLower(defaultAction)
	if defaultActionSmallCase != "y" && defaultActionSmallCase != "n" {
		return res, fmt.Errorf("Expecting values(Small&Upper case): \"y\"/\"n\"/\"\". But %s given", defaultAction)
	}
	var optionsStr string
	if defaultActionSmallCase == "y" {
		optionsStr = "Y/n"
	} else {
		optionsStr = "y/N"
	}
	prompt := promptui.Prompt{
		Label: fmt.Sprintf("%s [%s]", question, optionsStr),
		Validate: func(input string) error {
			inputSmallCase := strings.ToLower(input)
			if inputSmallCase == "y" || inputSmallCase == "n" || input == "" {
				return nil
			}
			return errors.New(`Prompt only accepts following values(Small&Upper case): "y"/"n"/""`)
		},
		Default: defaultActionSmallCase,
	}
	response, err := prompt.Run()
	if err != nil {
		return res, err
	}
	res = strings.ToLower(response) == "y" || response == ""
	return res, nil
}
