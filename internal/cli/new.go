package cli

import (
	"fmt"
	"time"

	"github.com/elhamza90/lifelog/internal/cli/io"
	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/elhamza90/lifelog/internal/http/rest/client"
	"github.com/elhamza90/lifelog/internal/http/rest/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new entity",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var newTagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Create a new Tag",
	Run: func(cmd *cobra.Command, args []string) {
		var tag domain.Tag
		if err := io.TagPrompt(&tag); err != nil {
			fmt.Println(err)
			return
		}
		token := viper.Get("Access").(string)
		id, err := client.PostTag(tag, token)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("\tSuccess: \n\tID: %d\n", id)
	},
}

var newActivityCmd = &cobra.Command{
	Use:   "activity",
	Short: "Create a new Activity",
	Run: func(cmd *cobra.Command, args []string) {
		token := viper.Get("Access").(string)
		tags, err := client.FetchTags(token)
		if err != nil {
			fmt.Println(err)
			return
		}
		activity := server.JSONReqActivity{
			Label:    "",
			Place:    "",
			Desc:     "",
			Time:     time.Now().Add(time.Duration(-1 * time.Hour)),
			Duration: time.Duration(time.Hour - (time.Minute * 15)),
		}
		if err := io.ActivityPrompt(&activity, tags); err != nil {
			fmt.Println(err)
			return
		}
		id, err := client.PostActivity(activity, token)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("\tSuccess: \n\tID: %d\n", id)
	},
}

var newExpenseCmd = &cobra.Command{
	Use:   "expense",
	Short: "Create a new Expense",
	Run: func(cmd *cobra.Command, args []string) {
		token := viper.Get("Access").(string)
		tags, err := client.FetchTags(token)
		if err != nil {
			fmt.Println(err)
			return
		}
		activities, err := client.FetchActivities(token, time.Now().AddDate(0, -3, 0))
		if err != nil {
			fmt.Println(err)
			return
		}
		expense := domain.Expense{
			Label: "",
			Value: 0,
			Unit:  "dhs",
			Time:  time.Now().Add(time.Duration(-1 * time.Hour)),
		}
		if err := io.ExpensePrompt(&expense, tags, activities); err != nil {
			fmt.Println(err)
			return
		}
		id, err := client.PostExpense(expense, token)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("\tSuccess: \n\tID: %d\n", id)
	},
}

func init() {
	newCmd.AddCommand(newTagCmd)
	newCmd.AddCommand(newActivityCmd)
	newCmd.AddCommand(newExpenseCmd)
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
