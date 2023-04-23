package cmd

import (
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/tordsk/sebastian/internal/sebastian"
	"os"
)

func init() {
	rootCmd.AddCommand(recipeCmd)
}

func selectRecipe(s *sebastian.Sebastian, cmd *cobra.Command) (name string) {
	if recipeName != "" {
		return recipeName
	}
	recipes, err := s.ListRecipes(cmd.Context())
	if err != nil {
		cmd.PrintErrln(err)
	}
	if len(recipes) == 0 {
		cmd.Help()
		return
	}
	recipes = append(recipes, "*Create new recipe*")
	q := promptui.Select{
		Label: "Select a recipe",
		Items: recipes,
	}
	i, name, err := q.Run()
	if err != nil {
		cmd.PrintErrln(err)
	}
	if i == len(recipes)-1 {
		recipeCreateCmd.Run(cmd, []string{})
		return selectRecipe(s, cmd)
	}
	return name
}

var recipeCmd = &cobra.Command{
	Use:   "recipe",
	Short: "list and manage recipes",
	Run:   interactive,
}

func interactive(cmd *cobra.Command, args []string) {
	s := sebastian.NewSebastian()
	recipeName = selectRecipe(s, cmd)

	recipeGetCmd.Run(cmd, []string{})

	intent := promptui.Select{
		Label: "What do you want to do?",
		Items: []string{
			"Comment",
			"Adjust",
			"Apply adjustments",
			"Delete",
			"Exit",
		},
	}
	_, intentResult, err := intent.Run()
	if err != nil {
		cmd.PrintErrln(err)
		os.Exit(0)
	}

	switch intentResult {
	case "Adjust":
		recipeAdjustCmd.Run(cmd, []string{})
	case "Apply adjustments":
		recipeAdjustApplyCmd.Run(cmd, []string{})
	case "Comment":
		recipeCommentCmd.Run(cmd, []string{})
	case "Delete":
		recipeDeleteCmd.Run(cmd, []string{})
	case "Exit":
		return
	}

	interactive(cmd, []string{})
}
