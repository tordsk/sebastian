package storage

import "time"

type RecipeData struct {
	Recipe      Recipe
	Adjustments []Adjustment
	Comments    []Comment
}

type Adjustment struct {
	Content    string
	Considered bool
}
type Comment struct {
	Content string
	Date    time.Time
}

func (c Comment) String() string {
	return c.Date.Format("2006-01-02") + ": " + c.Content
}

type Recipe string
type RecipeName string

type RecipeStorage interface {
	// GetRecipe returns a recipe by name
	GetRecipe(name RecipeName) (RecipeData, error)
	// GetRecipes returns all recipes
	GetRecipes() ([]RecipeName, error)
	// CreateRecipe creates a new recipe
	CreateRecipe(name RecipeName, recipe Recipe) error
	// UpdateRecipe updates an existing recipe
	ApplyAdjustments(name RecipeName, recipe Recipe) error

	CreateAdjustment(name RecipeName, result string) error
	DeleteAdjustment(name RecipeName, i int) error

	CreateComment(name RecipeName, comment string) error
	DeleteComment(name RecipeName, i int) error

	// DeleteRecipe deletes a recipe by name
	DeleteRecipe(name RecipeName) error
}
