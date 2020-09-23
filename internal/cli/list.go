package cli

import (
	"fmt"

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
		token := viper.Get("Access").(string)
		tags, err := client.FetchTags(token)
		if err != nil {
			fmt.Println(err)
			return
		}
		selectedTagIndex, err := tagSelect(tags)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Selected Tag: %v\n", tags[selectedTagIndex])
		return
	},
}

var listActivitiesCmd = &cobra.Command{
	Use:   "activities",
	Short: "List All Activities",
	Run: func(cmd *cobra.Command, args []string) {
		token := viper.Get("Access").(string)
		activities, err := client.FetchActivities(token)
		if err != nil {
			fmt.Println(err)
			return
		}
		selectedIndex, err := activitySelect(activities)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Selected Activities: %v\n", activities[selectedIndex])
		return
	},
}

func init() {
	listCmd.AddCommand(listTagCmd)
	listCmd.AddCommand(listActivitiesCmd)
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
