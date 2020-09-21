package cli

import (
	"fmt"

	"github.com/elhamza90/lifelog/internal/domain"
	"github.com/elhamza90/lifelog/internal/http/rest/client"
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
		tagToCreate, err := tagPrompt(domain.Tag{})
		if err != nil {
			fmt.Println(err)
			return
		}
		token := viper.Get("Access").(string)
		id, err := client.PostTag(tagToCreate, token)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("\tSuccess: \n\tID: %d\n", id)
	},
}

func init() {
	newCmd.AddCommand(newTagCmd)
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
