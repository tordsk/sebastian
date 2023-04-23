package cmd

import (
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/tordsk/sebastian/internal/sebastian"
)

func init() {

	recipeCmd.AddCommand(recipeCommentCmd)
}

var recipeCommentCmd = &cobra.Command{
	Use:   "comment",
	Short: "comment on how the recipe went ",
	Run: func(cmd *cobra.Command, args []string) {
		s := sebastian.NewSebastian()
		recipeName = selectRecipe(s, cmd)

		inputCommentPrompt := promptui.Prompt{
			Label: "Any comments on how the recipe went?",
		}
		comment, err := inputCommentPrompt.Run()
		if err != nil {
			cmd.PrintErrln(err)
		}

		err = s.CreateComment(cmd.Context(), recipeName, comment)
		if err != nil {
			cmd.PrintErrln(err)
		}
		cmd.Println("Comment saved")
	},
}
