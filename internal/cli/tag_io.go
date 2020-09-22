package cli

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"text/template"

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
	var idMaxLen int = 0
	for _, t := range tags {
		idStr := strconv.Itoa(int(t.ID))
		if len(idStr) > idMaxLen {
			idMaxLen = len(idStr)
		}
	}
	log.Println(idMaxLen)
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Inactive: "\t - " + `{{ printf "[%d] " .ID}} ` + fmt.Sprintf("{{ fixSpaces .ID %d}}", idMaxLen) + ` | {{ .Name | white }}`,
		Active:   "\t → " + `{{ printf "[%d]  " .ID|cyan|bold}} ` + fmt.Sprintf("{{ fixSpaces .ID %d}}", idMaxLen) + ` | {{ .Name | cyan | bold }}`,
		Selected: "\t → " + `{{ printf "[%d] " .ID| green | bold }}` + ` | {{ .Name | green | bold }}`,
		FuncMap: template.FuncMap{
			"fixSpaces": func(id domain.TagID, maxLen int) string {
				times := maxLen - len(strconv.Itoa(int(id)))
				return strings.Repeat(" ", times)
			},
			"black":     promptui.Styler(promptui.FGBlack),
			"red":       promptui.Styler(promptui.FGRed),
			"green":     promptui.Styler(promptui.FGGreen),
			"yellow":    promptui.Styler(promptui.FGYellow),
			"blue":      promptui.Styler(promptui.FGBlue),
			"magenta":   promptui.Styler(promptui.FGMagenta),
			"cyan":      promptui.Styler(promptui.FGCyan),
			"white":     promptui.Styler(promptui.FGWhite),
			"bgBlack":   promptui.Styler(promptui.BGBlack),
			"bgRed":     promptui.Styler(promptui.BGRed),
			"bgGreen":   promptui.Styler(promptui.BGGreen),
			"bgYellow":  promptui.Styler(promptui.BGYellow),
			"bgBlue":    promptui.Styler(promptui.BGBlue),
			"bgMagenta": promptui.Styler(promptui.BGMagenta),
			"bgCyan":    promptui.Styler(promptui.BGCyan),
			"bgWhite":   promptui.Styler(promptui.BGWhite),
			"bold":      promptui.Styler(promptui.FGBold),
			"faint":     promptui.Styler(promptui.FGFaint),
			"italic":    promptui.Styler(promptui.FGItalic),
			"underline": promptui.Styler(promptui.FGUnderline),
		},
	}
	tagPrompt := promptui.Select{
		Label:     "Choose Tag:",
		Items:     tags,
		Templates: templates,
		Size:      len(tags),
		Searcher: func(input string, index int) bool {
			name := strings.ToLower((tags)[index].Name)
			input = strings.ToLower(input)
			return strings.Contains(name, input)
		},
	}
	selectedTagIndex, _, err = tagPrompt.Run()
	return selectedTagIndex, err
}
