package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/elhamza90/lifelog/internal/http/rest/client"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lifelog",
	Short: "Lifelog is an Application to manage life activities",
	// Always check if Access token and config file exist before every command.
	// If necessary authenticate.
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if viper.ConfigFileUsed() == "" {
			fmt.Println("No Config File Found !")
			return
		}
		existingAccessToken := viper.Get("Access")
		// If Access Token was fetched and saved previously don't login
		if existingAccessToken == nil {
			fmt.Println("Logging in ...\n")
			pass, err := loginPrompt()
			if err != nil {
				fmt.Println(err)
				return
			}
			access, refresh, err := client.Login(pass)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Authentication Successful !\n")
			fmt.Println("Saving Token Pair ...\n")
			viper.Set("Access", access)
			viper.Set("Refresh", refresh)
			viper.WriteConfig()
			return
		}
		token, err := jwt.Parse(existingAccessToken.(string), nil)
		if token == nil {
			fmt.Println(err)
			return
		}
		claims, _ := token.Claims.(jwt.MapClaims)
		if expTime := time.Unix(int64(claims["exp"].(float64)), 0); !expTime.After(time.Now()) {
			refreshToken := viper.Get("Refresh").(string)
			newAccessToken, err := client.RefreshToken(refreshToken)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Refresh Token Successful")
			viper.Set("Access", newAccessToken)
			viper.WriteConfig()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lifelog.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".lifelog" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".lifelog")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
