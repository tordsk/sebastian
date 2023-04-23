package storage

import "time"

type InMemory struct {
	recipes map[RecipeName]RecipeData
}

func (im *InMemory) GetRecipe(name RecipeName) (RecipeData, error) {
	return im.recipes[name], nil
}

func (im *InMemory) GetRecipes() ([]RecipeName, error) {
	recipes := []RecipeName{}
	for name, _ := range im.recipes {
		recipes = append(recipes, name)
	}
	return recipes, nil
}

func (im *InMemory) CreateRecipe(name RecipeName, recipe Recipe) error {
	im.recipes[name] = RecipeData{Recipe: recipe}
	return nil
}

func (im *InMemory) ApplyAdjustments(name RecipeName, recipe Recipe) error {
	r := im.recipes[name]
	r.Recipe = recipe
	for i := 0; i < len(r.Adjustments); i++ {
		r.Adjustments[i].Considered = true
	}
	im.recipes[name] = r
	return nil
}

func (im *InMemory) DeleteRecipe(name RecipeName) error {
	delete(im.recipes, name)
	return nil
}

func (im *InMemory) CreateAdjustment(name RecipeName, result string) error {
	recipe, err := im.GetRecipe(name)
	if err != nil {
		return err
	}
	recipe.Adjustments = append(im.recipes[name].Adjustments, Adjustment{Content: result})
	im.recipes[name] = recipe
	return nil
}

func (im *InMemory) DeleteAdjustment(name RecipeName, i int) error {
	recipe, err := im.GetRecipe(name)
	if err != nil {
		return err
	}
	recipe.Adjustments = append(recipe.Adjustments[:i], recipe.Adjustments[i+1:]...)
	im.recipes[name] = recipe
	return nil
}

func (im *InMemory) CreateComment(name RecipeName, result string) error {
	recipe, err := im.GetRecipe(name)
	if err != nil {
		return err
	}
	recipe.Comments = append(im.recipes[name].Comments, Comment{Content: result, Date: time.Now()})
	im.recipes[name] = recipe
	return nil
}

func (im *InMemory) DeleteComment(name RecipeName, i int) error {
	recipe, err := im.GetRecipe(name)
	if err != nil {
		return err
	}
	recipe.Comments = append(recipe.Comments[:i], recipe.Comments[i+1:]...)
	im.recipes[name] = recipe
	return nil
}

func NewInMemory() *InMemory {
	return &InMemory{
		recipes: make(map[RecipeName]RecipeData),
	}
}
