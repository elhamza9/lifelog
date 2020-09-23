package cli

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/manifoldco/promptui"
)

func expensePrompt(defaultExpense domain.Expense, tags []domain.Tag, activities []domain.Activity) (t domain.Expense, err error) {
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
		selectedActivityIndex, err := activitySelect(activities)
		if err != nil {
			return t, err
		}
		activityID = activities[selectedActivityIndex].ID
	}
	log.Printf("Selected Activity ID: %d\n", activityID)
	// Label
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
	exp := domain.Expense{
		Label:      name,
		Value:      float32(value),
		Unit:       unit,
		Time:       time,
		ActivityID: activityID,
		Tags:       selectedTags,
	}
	return exp, nil
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

// expenseSelect list given expenses and asks user to select one.
func expenseSelect(expenses []domain.Expense) (selectedExpenseIndex int, err error) {
	var idMaxLen int = 0
	for _, act := range expenses {
		idStr := strconv.Itoa(int(act.ID))
		if len(idStr) > idMaxLen {
			idMaxLen = len(idStr)
		}
	}
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Inactive: "\t → " + `{{ printf "[%d]  " .ID}} ` + fmt.Sprintf("{{ fixSpaces .ID %d}}", idMaxLen) + ` | {{ .Time.Format "Mon Jan 02 2006" | white }}` + ` | {{ .Label | white }}` + ` {{ printf "(%.2f %s )" .Value .Unit | white}}`,
		Active:   "\t → " + `{{ printf "[%d]  " .ID}} ` + fmt.Sprintf("{{ fixSpaces .ID %d}}", idMaxLen) + ` | {{ .Time.Format "Mon Jan 02 2006" | cyan | bold }}` + ` | {{ .Label | cyan | bold }}` + ` {{ printf "(%.2f %s )" .Value .Unit | cyan | bold}}`,
		Selected: "\t → " + `{{ printf "[%d]  " .ID}} ` + fmt.Sprintf("{{ fixSpaces .ID %d}}", idMaxLen) + ` | {{ .Time.Format "Mon Jan 02 2006" | green | bold }}` + ` | {{ .Label | green | bold }}` + ` {{ printf "(%.2f %s )" .Value .Unit | green | bold}}`,
		FuncMap: template.FuncMap{
			"fixSpaces": func(id domain.ExpenseID, maxLen int) string {
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
	expensePrompt := promptui.Select{
		Label:     "Choose Expense:",
		Items:     expenses,
		Templates: templates,
		Size:      len(expenses),
		Searcher: func(input string, index int) bool {
			label := strings.ToLower(expenses[index].Label)
			input = strings.ToLower(input)
			return strings.Contains(label, input)
		},
	}
	selectedExpenseIndex, _, err = expensePrompt.Run()
	return selectedExpenseIndex, err
}
