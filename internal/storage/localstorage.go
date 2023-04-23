package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

type FileStorage struct {
	file     string
	InMemory *InMemory
}

func (f *FileStorage) load() error {
	b, err := os.ReadFile(f.file)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &f.InMemory.recipes)
}

func (f *FileStorage) persist() error {
	b, err := json.Marshal(f.InMemory.recipes)
	if err != nil {
		return err
	}
	return os.WriteFile(f.file, b, 0644)
}

func (f *FileStorage) GetRecipe(name RecipeName) (RecipeData, error) {
	f.load()
	return f.InMemory.GetRecipe(name)
}

func (f *FileStorage) GetRecipes() ([]RecipeName, error) {
	f.load()
	return f.InMemory.GetRecipes()
}

func (f *FileStorage) CreateRecipe(name RecipeName, recipe Recipe) error {
	f.load()
	defer f.persist()
	return f.InMemory.CreateRecipe(name, recipe)
}

func (f *FileStorage) ApplyAdjustments(name RecipeName, recipe Recipe) error {
	f.load()
	defer f.persist()
	return f.InMemory.ApplyAdjustments(name, recipe)
}

func (f *FileStorage) DeleteRecipe(name RecipeName) error {
	f.load()
	defer f.persist()
	return f.InMemory.DeleteRecipe(name)
}

func (f *FileStorage) CreateAdjustment(name RecipeName, result string) error {
	f.load()
	defer f.persist()
	return f.InMemory.CreateAdjustment(name, result)
}

func (f *FileStorage) DeleteAdjustment(name RecipeName, i int) error {
	f.load()
	defer f.persist()
	return f.InMemory.DeleteAdjustment(name, i)
}

func (f *FileStorage) CreateComment(name RecipeName, comment string) error {
	f.load()
	defer f.persist()
	return f.InMemory.CreateComment(name, comment)
}

func (f *FileStorage) DeleteComment(name RecipeName, i int) error {
	f.load()
	defer f.persist()
	return f.InMemory.DeleteComment(name, i)
}

func NewFileStorage() *FileStorage {
	dirname, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	configDir := fmt.Sprintf("%s/.sebastian", dirname)
	configFile := fmt.Sprintf("%s/data.json", configDir)
	err = os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return &FileStorage{
		file:     configFile,
		InMemory: NewInMemory(),
	}
}
