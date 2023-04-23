package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tordsk/sebastian/internal/sebastian"
)

func init() {
	recipeCmd.AddCommand(recipeDeleteCmd)
}

var recipeDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a recipe",
	Run: func(cmd *cobra.Command, args []string) {
		s := sebastian.NewSebastian()
		recipeName = selectRecipe(s, cmd)

		err := s.DeleteRecipe(cmd.Context(), recipeName)
		if err != nil {
			cmd.PrintErrln(err)
		}
		recipeName = ""
		cmd.Println("Deleted")
	},
}
