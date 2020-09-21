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

// tagSelect list given tags and asks user to select one.
func tagSelect(tags []domain.Tag) (selectedTagIndex int, err error) {
	/*
		var idMaxLen int = 0
		for _, t := range tags {
			idStr := strconv.Itoa(int(t.ID))
			if len(idStr) > idMaxLen {
				idMaxLen = len(idStr)
			}
		}
	*/
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Inactive: "\t - " + `{{ printf "[%d] " .ID | white }}` + `| {{ .Name | white }}`,
		Active: "\t → " +
			`{{ printf "[%d] " .ID | cyan | bold }}` +
			`| {{ .Name | cyan | bold }}`,
		Selected: "\t → " +
			`{{ printf "[%d] " .ID | green | bold }}` +
			`| {{ .Name | green | bold }}`,
	}
	tagPrompt := promptui.Select{
		Label:     "Choose Tag:",
		Items:     tags,
		Templates: templates,
		Size:      len(tags),
		/*Searcher: func(input string, index int) bool {
			name := strings.ToLower((tags)[index].Name)
			input = strings.ToLower(input)
			return strings.Contains(name, input)
		},*/
	}
	selectedTagIndex, _, err = tagPrompt.Run()
	return selectedTagIndex, err
}
