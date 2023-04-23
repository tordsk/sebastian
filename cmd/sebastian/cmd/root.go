package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	// Flag vars
	recipeName string
)

var rootCmd = &cobra.Command{
	Use:   "sebastian",
	Short: "sebastian is a great chef and a great friend",
	Long: `Sebastian will help you cook your favorite recipes.
Adjust the ingredients and the cooking time to your needs. 
And replace the ingredients with your own favorites.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
