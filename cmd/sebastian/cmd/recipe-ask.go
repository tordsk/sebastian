package cmd

import (
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/tordsk/sebastian/internal/sebastian"
)

func init() {
	recipeCmd.AddCommand(recipeAskCmd)
}

var recipeAskCmd = &cobra.Command{
	Use:   "ask",
	Short: "ask a question or get some help during cooking",
	Run: func(cmd *cobra.Command, args []string) {
		s := sebastian.NewSebastian()
		recipeName = selectRecipe(s, cmd)

		inputCommentPrompt := promptui.Prompt{
			Label: "What would you like to ask the chef?",
		}
		comment, err := inputCommentPrompt.Run()
		if err != nil {
			cmd.PrintErrln(err)
		}

		outStream, err := s.AskTheChef(cmd.Context(), recipeName, comment)
		if err != nil {
			cmd.PrintErrln(err)
		}
		sebastian.LogAndGetString(cmd, outStream)
	},
}
