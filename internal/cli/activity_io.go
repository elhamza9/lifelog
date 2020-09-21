package cli

import (
	"errors"
	"fmt"
	"time"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/manifoldco/promptui"
)

// activityInput gets activity fields from user and returns a new activity with filled fields
func activityPrompt(defaultActivity domain.Activity, tags []domain.Tag) (a domain.Activity, err error) {
	// Label Prompt
	prompt := promptui.Prompt{
		Label:    "Label",
		Validate: activityLabelValidator,
		Default:  defaultActivity.Label,
	}
	label, err := prompt.Run()
	if err != nil {
		return a, err
	}
	// Place Prompt
	prompt = promptui.Prompt{
		Label:    "Place",
		Validate: activityPlaceValidator,
		Default:  defaultActivity.Place,
	}
	place, err := prompt.Run()
	if err != nil {
		return a, err
	}
	// Description Prompt
	prompt = promptui.Prompt{
		Label:    "Description",
		Validate: activityDescValidator,
		Default:  defaultActivity.Desc,
	}
	desc, err := prompt.Run()
	if err != nil {
		return a, err
	}
	// Time Start Prompt
	prompt = promptui.Prompt{
		Label:    "Time Start",
		Validate: activityTimeValidator,
		Default:  defaultActivity.Time.Format("2006-01-02 15:04"),
	}
	actTimeStr, err := prompt.Run()
	if err != nil {
		return a, err
	}
	zone, _ := time.Now().Zone()
	actTime, err := time.Parse("2006-01-02 15:04 MST", actTimeStr+" "+zone)
	if err != nil {
		return a, err
	}
	// Duration Prompt
	prompt = promptui.Prompt{
		Label:    "Duration",
		Validate: activityDurationValidator,
		Default:  defaultActivity.Duration.String(),
	}
	durationStr, err := prompt.Run()
	if err != nil {
		return a, err
	}
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return a, err
	}
	// Tags
	noTag := domain.Tag{ID: 0, Name: "OK"}
	tags = append(tags, noTag)
	selectedTags := []domain.Tag{}
	// Run infinite loop. Break when Tag noTag is selected
	for {
		selectedTagIndex, err := tagSelect(tags)
		if err != nil {
			return a, err
		}
		selectedTag := tags[selectedTagIndex]
		if selectedTag.ID == noTag.ID {
			break
		} else {
			selectedTags = append(selectedTags, selectedTag)
			// Remove selected Tag from list
			tags = append(tags[:selectedTagIndex], tags[selectedTagIndex+1:]...)
		}
	}
	a = domain.Activity{
		Label:    label,
		Place:    place,
		Desc:     desc,
		Time:     actTime,
		Duration: duration,
		Tags:     selectedTags,
	}
	return a, nil
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
