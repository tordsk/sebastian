package cmd

import (
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/tordsk/sebastian/internal/sebastian"
	"github.com/tordsk/sebastian/internal/storage"
)

func init() {
	recipeCmd.AddCommand(recipeCreateCmd)
}

var recipeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create recipes",
	Run: func(cmd *cobra.Command, args []string) {
		s := sebastian.NewSebastian()

		whatDish := promptui.Prompt{
			Label: "What dish should I create a recipe for?",
		}
		recipeName, err := whatDish.Run()

		recipeChan, err := s.NewRecipe(cmd.Context(), recipeName)
		if err != nil {
			cmd.PrintErrln(err)
		}

		fullRecipe := sebastian.LogAndGetString(cmd, recipeChan)

		approve := promptui.Select{
			Label: "Keep this recipe?",
			Items: []string{"Yes", "No"},
		}
		_, save, err := approve.Run()

		if err != nil {
			cmd.PrintErrln(err)
		}
		if save == "No" {
			cmd.Println("Discarding recipe, try again with a different input")
			return
		}
		recipeName, err = s.CreateRecipe(cmd.Context(), recipeName, storage.Recipe(fullRecipe))
		cmd.Println("Saved!")
	},
}
