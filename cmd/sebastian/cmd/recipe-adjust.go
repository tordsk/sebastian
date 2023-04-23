package cmd

import (
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/tordsk/sebastian/internal/sebastian"
)

func init() {

	recipeCmd.AddCommand(recipeAdjustCmd)
}

var recipeAdjustCmd = &cobra.Command{
	Use:   "adjust",
	Short: "add some info for how to adjust the recipe ",
	Run: func(cmd *cobra.Command, args []string) {
		s := sebastian.NewSebastian()
		recipeName = selectRecipe(s, cmd)

		inputCommentPrompt := promptui.Prompt{
			Label: "Any comments to consider when adjusting the recipe?",
		}
		comment, err := inputCommentPrompt.Run()
		if err != nil {
			cmd.PrintErrln(err)
		}

		err = s.CreateAdjustment(cmd.Context(), recipeName, comment)
		if err != nil {
			cmd.PrintErrln(err)
		}
		cmd.Println("Result saved")
	},
}
