package cli

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/manifoldco/promptui"
)

func expensePrompt(defaultExpense domain.Expense, tags []domain.Tag) (t domain.Expense, err error) {
	/*
		// Activity
		var yesNoPromptQuestion string
		if defaultExpense.ActivityID > 0 {
			yesNoPromptQuestion = fmt.Sprintf("Change Activity [%d] ?", defaultExpense.ActivityID)
		} else {
			yesNoPromptQuestion = "Select an Activity ?"
		}
		selectActivity, err := yesNoPrompt(yesNoPromptQuestion, "N")
		if err != nil {
			fmt.Println(err)
			return
		}
		activityID := defaultExpense.ActivityID
		if selectActivity {
			activities, err := store.FetchActivities(time.Now().AddDate(0, -3, 0).Format("2006-01-02"))
			if err != nil {
				return t, err
			}
			selectedActivityIndex, err := activitySelect(&activities)
			if err != nil {
				return t, err
			}
			activityID = activities[selectedActivityIndex].ID
		}
	*/
	prompt := promptui.Prompt{
		Label:    "Label",
		Validate: expenseLabelValidator,
		Default:  defaultExpense.Label,
	}
	name, err := prompt.Run()
	if err != nil {
		return t, err
	}
	// Value
	prompt = promptui.Prompt{
		Label:    "Value",
		Validate: expenseValueValidator,
		Default:  fmt.Sprintf("%.2f", defaultExpense.Value),
	}
	valueStr, err := prompt.Run()
	if err != nil {
		return t, err
	}
	value, _ := strconv.ParseFloat(valueStr, 32)
	// Unit
	prompt = promptui.Prompt{
		Label:    "Unit",
		Validate: expenseUnitValidator,
		Default:  defaultExpense.Unit,
	}
	unit, err := prompt.Run()
	if err != nil {
		return t, err
	}
	// Time
	if defaultExpense.Time.IsZero() {
		defaultExpense.Time = time.Now()
	}
	prompt = promptui.Prompt{
		Label:    "Time",
		Validate: expenseTimeValidator,
		Default:  defaultExpense.Time.Format("2006-01-02"),
	}
	timeStr, err := prompt.Run()
	if err != nil {
		return t, err
	}
	time, _ := time.Parse("2006-01-02", timeStr)
	// Tags
	noTag := domain.Tag{ID: 0, Name: "OK"}
	tags = append(tags, noTag)
	selectedTags := []domain.Tag{}
	// Run infinite loop. Break when Tag noTag is selected
	for {
		selectedTagIndex, err := tagSelect(tags)
		if err != nil {
			return domain.Expense{}, err
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
	t = domain.Expense{
		Label: name,
		Value: float32(value),
		Unit:  unit,
		Time:  time,
		//		ActivityID: activityID,
		Tags: selectedTags,
	}
	return t, nil
}

// expenseLabelValidator validates the expense name inputed by the user
func expenseLabelValidator(input string) error {
	if len(input) < 2 {
		return errors.New("Label must have at least 2 characters")
	} else if len(input) > 50 {
		return errors.New("Label must have less than 50 characters")
	} else {
		return nil
	}
}

// expenseValueValidator validates the expense value inputed by the user
func expenseValueValidator(input string) error {
	res, err := strconv.ParseFloat(input, 32)
	if err == nil && res <= 0 {
		return errors.New("Value must be greater than or equal to zero")
	}
	return err
}

// expenseUnitValidator validates the expense unit inputed by the user
func expenseUnitValidator(input string) error {
	if len(input) == 0 {
		return errors.New("Unit must have at least 1 characters")
	} else if len(input) > 50 {
		return errors.New("Unit must have less than 50 characters")
	} else {
		return nil
	}
}

// expenseTimeValidator validates the expense time inputed by the user
func expenseTimeValidator(input string) error {
	res, err := time.Parse("2006-01-02", input)
	if err == nil && res.After(time.Now()) {
		return errors.New("Expense Time can not be future")
	}
	return err
}
