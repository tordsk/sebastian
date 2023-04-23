package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tordsk/sebastian/internal/sebastian"
	"github.com/tordsk/sebastian/internal/storage"
)

func init() {
	recipeAdjustCmd.AddCommand(recipeAdjustApplyCmd)
}

var recipeAdjustApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "apply previous comments to adjust the recipe",
	Run: func(cmd *cobra.Command, args []string) {
		s := sebastian.NewSebastian()
		recipeName = selectRecipe(s, cmd)

		adjustments, err := s.ApplyAdjustments(cmd.Context(), recipeName)
		if err != nil {
			return
		}
		fullRecipe := sebastian.LogAndGetString(cmd, adjustments)
		s.AdjustRecipe(cmd.Context(), recipeName, storage.Recipe(fullRecipe))
	},
}
