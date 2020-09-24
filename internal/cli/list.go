package cli

import (
	"fmt"
	"strconv"
	"time"

	"github.com/elhamza90/lifelog/internal/cli/io"
	"github.com/elhamza90/lifelog/internal/http/rest/client"
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
		tag := tags[selectedTagIndex]
		// Action
		action, err := io.ActionPrompt()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch action {
		case io.ActionDelete:
			if err := client.DeleteTag(tag.ID, token); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Tag deleted successfully")
		case io.ActionEdit:
			if err := io.TagPrompt(&tag); err != nil {
				fmt.Println(err)
				return
			}
			if err := client.UpdateTag(tag, token); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Tag updated successfully")
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
		activity := activities[selectedIndex]
		// Action
		action, err := io.ActionPrompt()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch action {
		case io.ActionDelete:
			if err := client.DeleteActivity(activity.ID, token); err != nil {
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
		expense := expenses[selectedIndex]
		switch action {
		case io.ActionDelete:
			if err := client.DeleteExpense(expense.ID, token); err != nil {
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
			res, err := client.FetchExpenseDetails(expense.ID, token)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(res)
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
