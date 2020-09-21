package cli

import (
	"errors"
	"strings"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/manifoldco/promptui"
)

func tagPrompt(defaultTag domain.Tag) (t domain.Tag, err error) {
	prompt := promptui.Prompt{
		Label:    "Name",
		Validate: tagNameValidator,
		Default:  defaultTag.Name,
	}
	name, err := prompt.Run()
	if err != nil {
		return t, err
	}
	t = domain.Tag{
		Name: name,
	}
	return t, nil
}

// tagNameValidator validates the tag name inputed by the user
func tagNameValidator(input string) error {
	if len(input) < 2 {
		return errors.New("Tag Name must have at least 2 characters")
	} else if len(input) > 50 {
		return errors.New("Tag Name must have less than 50 characters")
	} else if strings.Contains(input, " ") || strings.Contains(input, "\t") {
		return errors.New("Tag name can not contain spaces")
	} else {
		return nil
	}
}
