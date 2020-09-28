package cli

import (
	"fmt"
	"strconv"
	"time"

	"github.com/elhamza90/lifelog/internal/cli/io"
	"github.com/elhamza90/lifelog/internal/http/rest/client"
	"github.com/elhamza90/lifelog/internal/http/rest/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List entities",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var listTagCmd = &cobra.Command{
	Use:   "tags",
	Short: "List All Tags",
	Run: func(cmd *cobra.Command, args []string) {
		// Fetch
		token := viper.Get("Access").(string)
		tags, err := client.FetchTags(token)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Select
		selectedTagIndex, err := io.TagSelect(tags)
		if err != nil {
			fmt.Println(err)
			return
		}
		selected := tags[selectedTagIndex]
		// Action
		action, err := io.ActionPrompt()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch action {
		case io.ActionDelete:
			if err := client.DeleteTag(selected.ID, token); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Tag deleted successfully")
		case io.ActionEdit:
			tag := server.JSONReqTag{
				ID:   selected.ID,
				Name: selected.Name,
			}
			if err := io.TagPrompt(&tag); err != nil {
				fmt.Println(err)
				return
			}
			if err := client.UpdateTag(tag, token); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Tag updated successfully")
		case io.ActionDetails:
			const (
				showExpenses   string = "Show Expenses"
				showActivities string = "Show Activities"
			)
			action, err := io.CustomPrompt("What do you want to do", []string{showActivities, showExpenses, "Exit"})
			if err != nil {
				fmt.Println(err)
				return
			}
			if action == showActivities {
				activities, err := client.FetchTagActivities(selected.ID, token)
				if err != nil {
					fmt.Println(err)
					return
				}
				if _, err = io.ActivitySelect(activities); err != nil {
					fmt.Println(err)
					return
				}
			} else if action == showExpenses {
				expenses, err := client.FetchTagExpenses(selected.ID, token)
				if err != nil {
					fmt.Println(err)
					return
				}
				if _, err = io.ExpenseSelect(expenses); err != nil {
					fmt.Println(err)
					return
				}
			} else {
				return
			}
		}
		return
	},
}

var activitiesMonthsFlag string
var listActivitiesCmd = &cobra.Command{
	Use:   "activities",
	Short: "List All Activities",
	Run: func(cmd *cobra.Command, args []string) {
		// Parse Flag
		nbrMonths, err := strconv.Atoi(activitiesMonthsFlag)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Fetch
		token := viper.Get("Access").(string)
		activities, err := client.FetchActivities(token, time.Now().AddDate(0, -nbrMonths, 0))
		if err != nil {
			fmt.Println(err)
			return
		}
		// Select
		selectedIndex, err := io.ActivitySelect(activities)
		if err != nil {
			fmt.Println(err)
			return
		}
		selected := activities[selectedIndex]
		// Action
		action, err := io.ActionPrompt()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch action {
		case io.ActionDelete:
			if err := client.DeleteActivity(selected.ID, token); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Activity deleted successfully")
		case io.ActionEdit:
			// Prompt
			tags, err := client.FetchTags(token)
			if err != nil {
				fmt.Println(err)
				return
			}
			activity := server.JSONReqActivity{
				ID:       selected.ID,
				Label:    selected.Label,
				Place:    selected.Place,
				Desc:     selected.Desc,
				Time:     selected.Time,
				Duration: selected.Duration,
			}
			if err := io.ActivityPrompt(&activity, tags); err != nil {
				fmt.Println(err)
				return
			}
			// Update
			if err := client.UpdateActivity(activity, token); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Activity updated successfully")
		case io.ActionDetails:
			res, err := client.FetchActivityDetails(selected.ID, token)
			if err != nil {
				fmt.Println(err)
				return
			}
			// Print
			fmt.Printf("Activity:\n\t- ID: %d\n\t- Label: %s\n\t- Time: %s\n\t- Place: %s\n\t- Duration: %s\n", res.ID, res.Label, res.Time.Format("2006-01-02 15:04"), res.Place, res.Duration)
			fmt.Printf("\t- Tags: ")
			for _, t := range res.Tags {
				fmt.Printf("- %s ", t.Name)
			}
			fmt.Println("\n\t- Expenses:")
			for _, exp := range res.Expenses {
				fmt.Printf("\t\t- %s (%.2f %s)\n", exp.Label, exp.Value, exp.Unit)
			}
			fmt.Println()
		}
		return
	},
}

var expensesMonthsFlag string
var listExpensesCmd = &cobra.Command{
	Use:   "expenses",
	Short: "List All Expenses",
	Run: func(cmd *cobra.Command, args []string) {
		// Parse Argument if exists
		nbrMonths, err := strconv.Atoi(expensesMonthsFlag)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Fetch
		token := viper.Get("Access").(string)
		expenses, err := client.FetchExpenses(token, time.Now().AddDate(0, -nbrMonths, 0))
		if err != nil {
			fmt.Println(err)
			return
		}
		// Select
		selectedIndex, err := io.ExpenseSelect(expenses)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Action
		action, err := io.ActionPrompt()
		if err != nil {
			fmt.Println(err)
			return
		}
		selected := expenses[selectedIndex]
		switch action {
		case io.ActionDelete:
			if err := client.DeleteExpense(selected.ID, token); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Expense deleted successfully")
		case io.ActionEdit:
			// Prompt
			tags, err := client.FetchTags(token)
			if err != nil {
				fmt.Println(err)
				return
			}
			activities, err := client.FetchActivities(token, time.Now().AddDate(0, -2, 0))
			if err != nil {
				fmt.Println(err)
				return
			}
			expense := server.JSONReqExpense{
				ID:         selected.ID,
				Label:      selected.Label,
				Time:       selected.Time,
				Value:      selected.Value,
				Unit:       selected.Unit,
				ActivityID: selected.ActivityID,
			}
			if err := io.ExpensePrompt(&expense, tags, activities); err != nil {
				fmt.Println(err)
				return
			}
			// Update
			if err := client.UpdateExpense(expense, token); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Expense updated successfully")
		case io.ActionDetails:
			res, err := client.FetchExpenseDetails(selected.ID, token)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("Expense:\n\t- ID: %d\n\t- Label: %s\n\t- Time: %s\n\t- Value: %.2f\n\t- Unit: %s\n\t- Activity ID: %d\n", res.ID, res.Label, res.Time.Format("2006-01-02 15:04"), res.Value, res.Unit, res.ActivityID)
			fmt.Printf("\t- Tags: ")
			for _, t := range res.Tags {
				fmt.Printf("- %s ", t.Name)
			}
			fmt.Println()
		}
		return
	},
}

func init() {
	listCmd.AddCommand(listTagCmd)
	listActivitiesCmd.Flags().StringVarP(&activitiesMonthsFlag, "months", "m", "3", "List activities occured in the last n months")
	listCmd.AddCommand(listActivitiesCmd)
	listExpensesCmd.Flags().StringVarP(&expensesMonthsFlag, "months", "m", "3", "List expenses occured in the last n months")
	listCmd.AddCommand(listExpensesCmd)
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
