package io

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/elhamza90/lifelog/internal/http/rest/server"
	"github.com/manifoldco/promptui"
)

// ActivityPrompt asks user to fill activity fields
func ActivityPrompt(activity *server.JSONReqActivity, tags []server.JSONRespListTag) error {
	// Label Prompt
	prompt := promptui.Prompt{
		Label:    "Label",
		Validate: activityLabelValidator,
		Default:  (*activity).Label,
	}
	label, err := prompt.Run()
	if err != nil {
		return err
	}
	// Place Prompt
	prompt = promptui.Prompt{
		Label:    "Place",
		Validate: activityPlaceValidator,
		Default:  (*activity).Place,
	}
	place, err := prompt.Run()
	if err != nil {
		return err
	}
	// Description Prompt
	prompt = promptui.Prompt{
		Label:    "Description",
		Validate: activityDescValidator,
		Default:  (*activity).Desc,
	}
	desc, err := prompt.Run()
	if err != nil {
		return err
	}
	// Time Start Prompt
	prompt = promptui.Prompt{
		Label:    "Time Start",
		Validate: activityTimeValidator,
		Default:  (*activity).Time.Format("2006-01-02 15:04"),
	}
	actTimeStr, err := prompt.Run()
	if err != nil {
		return err
	}
	zone, _ := time.Now().Zone()
	actTime, err := time.Parse("2006-01-02 15:04 MST", actTimeStr+" "+zone)
	if err != nil {
		return err
	}
	// Duration Prompt
	prompt = promptui.Prompt{
		Label:    "Duration",
		Validate: activityDurationValidator,
		Default:  (*activity).Duration.String(),
	}
	durationStr, err := prompt.Run()
	if err != nil {
		return err
	}
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return err
	}
	// Tags
	noTag := server.JSONRespListTag{ID: 0, Name: "OK"}
	tags = append(tags, noTag)
	selectedIds := []domain.TagID{}
	// Run infinite loop. Break when Tag noTag is selected
	for {
		selectedTagIndex, err := TagSelect(tags)
		if err != nil {
			return err
		}
		selectedTag := tags[selectedTagIndex]
		if selectedTag.ID == noTag.ID {
			break
		} else {
			selectedIds = append(selectedIds, selectedTag.ID)
			// Remove selected Tag from list
			tags = append(tags[:selectedTagIndex], tags[selectedTagIndex+1:]...)
		}
	}
	(*activity).Label = label
	(*activity).Place = place
	(*activity).Desc = desc
	(*activity).Time = actTime
	(*activity).Duration = duration
	(*activity).TagIds = selectedIds
	return nil
}

// Validators

// activityLabelValidator validates the activity label inputed by the user
func activityLabelValidator(input string) error {
	if len(input) < 3 {
		return errors.New("Label must have at least 3 characters")
	} else if len(input) > 50 {
		return errors.New("Label must have less than 50 characters")
	} else {
		return nil
	}
}

// activityPlaceValidator validates the activity place inputed by the user
func activityPlaceValidator(input string) error {
	if len(input) > 50 {
		return errors.New("Place must have less than 50 characters")
	}
	return nil
}

// activityDescValidator validates the activity description inputed by the user
func activityDescValidator(input string) error {
	if len(input) > 300 {
		return errors.New("Description must have less than 300 characters")
	}
	return nil
}

// activityTimeValidator validates the activity time inputed by the user
func activityTimeValidator(input string) error {
	// get current time ignoring seconds
	now := time.Now().Round(time.Minute)
	zone, _ := now.Zone()
	input += " " + zone
	inputTime, err := time.Parse("2006-01-02 15:04 MST", input)
	if err != nil {
		return err
	}
	if inputTime.IsZero() {
		return errors.New("Input Time can not be Zero")
	}
	if inputTime.After(now) {
		return fmt.Errorf("Activity Time (%s) can not be future. Current Time is: %s. ", inputTime, now)
	}
	return nil
}

// activityDurationValidator validates the duration of activity inputed by user
func activityDurationValidator(input string) error {
	dur, err := time.ParseDuration(input)
	if err != nil {
		return err
	}
	if dur < 0 {
		return errors.New("Activity Duration can not be negative")
	}
	return nil
}

// ActivitySelect lists given activities and asks user to select one.
func ActivitySelect(activities []server.JSONRespListActivity) (selectedActivityIndex int, err error) {
	var idMaxLen int = 0
	for _, act := range activities {
		idStr := strconv.Itoa(int(act.ID))
		if len(idStr) > idMaxLen {
			idMaxLen = len(idStr)
		}
	}
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Inactive: "\t → " + `{{ printf "[%d]  " .ID}} ` + fmt.Sprintf("{{ fixSpaces .ID %d}}", idMaxLen) + ` | {{ .Time.Format "Mon Jan 02 2006" | white }}` + ` | {{ .Label | white }}`,
		Active:   "\t → " + `{{ printf "[%d]  " .ID}} ` + fmt.Sprintf("{{ fixSpaces .ID %d}}", idMaxLen) + ` | {{ .Time.Format "Mon Jan 02 2006" | cyan | bold }}` + ` | {{ .Label | cyan | bold }}`,
		Selected: "\t → " + `{{ printf "[%d]  " .ID}} ` + fmt.Sprintf("{{ fixSpaces .ID %d}}", idMaxLen) + ` | {{ .Time.Format "Mon Jan 02 2006" | green | bold }}` + ` | {{ .Label | green | bold }}`,
		FuncMap: template.FuncMap{
			"fixSpaces": func(id domain.ActivityID, maxLen int) string {
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
	activityPrompt := promptui.Select{
		Label:     "Choose Activity:",
		Items:     activities,
		Templates: templates,
		Size:      len(activities),
		Searcher: func(input string, index int) bool {
			label := strings.ToLower(activities[index].Label)
			input = strings.ToLower(input)
			return strings.Contains(label, input)
		},
	}
	selectedActivityIndex, _, err = activityPrompt.Run()
	return selectedActivityIndex, err
}
