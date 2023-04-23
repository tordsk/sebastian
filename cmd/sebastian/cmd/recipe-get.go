package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tordsk/sebastian/internal/sebastian"
	"sort"
)

func init() {
	recipeGetCmd.Flags().StringVarP(&recipeName, "name", "n", "", "name of recipe to get")
	recipeCmd.AddCommand(recipeGetCmd)
}

var recipeGetCmd = &cobra.Command{
	Use:   "get",
	Short: "list and view recipes",
	Run: func(cmd *cobra.Command, args []string) {
		s := sebastian.NewSebastian()
		recipeName = selectRecipe(s, cmd)

		recipe, err := s.GetRecipe(cmd.Context(), recipeName)
		if err != nil {
			return
		}
		sort.Slice(recipe.Comments, func(i, j int) bool {
			return recipe.Comments[i].Date.After(recipe.Comments[j].Date)
		})
		cmd.Println(recipe.Recipe)
		if len(recipe.Comments) > 0 {
			cmd.Println("Latest result: ")
			for i, comment := range recipe.Comments {
				if i > 4 {
					break
				}
				cmd.Println("-", comment)
			}
		}
	},
}
