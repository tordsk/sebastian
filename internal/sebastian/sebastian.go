package sebastian

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/tordsk/sebastian/internal/ai"
	"github.com/tordsk/sebastian/internal/prompts"
	"github.com/tordsk/sebastian/internal/storage"
	"io"
	"regexp"
	"strings"
)

const OrgIntangible = "org-HSrpBsfeNo5ESDY4zDesee7O"
const OrgKidBeer = "org-hvIJPl1LOKzqecl4yDguubtd"

type Sebastian struct {
	ai      *openai.Client
	storage storage.RecipeStorage
}

func NewSebastian() *Sebastian {
	s := &Sebastian{}

	if s.ai == nil {
		s.ai = ai.NewClient("", ai.WithOrgID(OrgIntangible))
	}

	if s.storage == nil {
		s.storage = storage.NewFileStorage()
	}

	return s
}

func (s *Sebastian) NewRecipe(ctx context.Context, name string) (chan string, error) {
	systemMsg := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: prompts.SystemPrompt,
	}

	prompt := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: fmt.Sprintf(prompts.CreateRecipe, name),
	}

	completion, err := s.ai.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:            openai.GPT3Dot5Turbo,
		Messages:         []openai.ChatCompletionMessage{systemMsg, prompt},
		MaxTokens:        2000,
		Temperature:      0,
		TopP:             1,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
	})
	if err != nil {
		return nil, err
	}

	return streamReceiver(completion), nil
}

func (s Sebastian) CreateRecipe(ctx context.Context, name string, recipe storage.Recipe) (string, error) {
	// override name if it's in the recipe
	re := regexp.MustCompile("Name:\\s+(.*)")
	match := re.FindStringSubmatch(string(recipe))
	if len(match) > 1 {
		name = match[1]
	}
	return name, s.storage.CreateRecipe(storage.RecipeName(name), recipe)
}

func (s Sebastian) ListRecipes(ctx context.Context) ([]storage.RecipeName, error) {
	return s.storage.GetRecipes()
}

func (s Sebastian) GetRecipe(ctx context.Context, name string) (storage.RecipeData, error) {
	return s.storage.GetRecipe(storage.RecipeName(name))
}

func (s Sebastian) DeleteRecipe(ctx context.Context, name string) error {
	return s.storage.DeleteRecipe(storage.RecipeName(name))
}

func (s Sebastian) AdjustRecipe(ctx context.Context, name string, recipe storage.Recipe) error {
	return s.storage.ApplyAdjustments(storage.RecipeName(name), recipe)
}

func (s Sebastian) ApplyAdjustments(ctx context.Context, name string) (chan string, error) {
	r, err := s.GetRecipe(ctx, name)
	if err != nil {
		return nil, err
	}
	systemMsg := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: prompts.AdjustSystemPrompt,
	}

	oldRecipe := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: string(r.Recipe),
	}

	var comments strings.Builder
	for _, result := range r.Adjustments {
		if !result.Considered {
			comments.WriteString(fmt.Sprintf("- %s\n", result.Content))
		}
	}
	applyAdjustments := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: fmt.Sprintf(prompts.AdjustRecipe, name, comments.String()),
	}

	completion, err := s.ai.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:            openai.GPT3Dot5Turbo,
		Messages:         []openai.ChatCompletionMessage{systemMsg, oldRecipe, applyAdjustments},
		MaxTokens:        2000,
		Temperature:      0,
		TopP:             1,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
	})
	if err != nil {
		return nil, err
	}

	return streamReceiver(completion), nil
}

func (s Sebastian) CreateAdjustment(ctx context.Context, name string, result string) error {
	return s.storage.CreateAdjustment(storage.RecipeName(name), result)
}

func (s Sebastian) AskTheChef(ctx context.Context, name string, question string) (chan string, error) {
	r, err := s.GetRecipe(ctx, name)
	if err != nil {
		return nil, err
	}
	systemMsg := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: prompts.AskTheChefSystemPrompt,
	}

	recipe := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: fmt.Sprintf("We will be making this recipe today: \n %s", string(r.Recipe)),
	}

	askTheChef := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: fmt.Sprintf(prompts.AskTheChef, question),
	}

	completion, err := s.ai.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:            openai.GPT3Dot5Turbo,
		Messages:         []openai.ChatCompletionMessage{systemMsg, recipe, askTheChef},
		MaxTokens:        2000,
		Temperature:      0,
		TopP:             1,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
	})
	if err != nil {
		return nil, err
	}

	return streamReceiver(completion), nil
}

func (s Sebastian) CreateComment(ctx context.Context, name string, comment string) error {
	return s.storage.CreateComment(storage.RecipeName(name), comment)
}

func streamReceiver(stream *openai.ChatCompletionStream) chan string {
	pipe := make(chan string)
	go func() {
		for {
			part, err := stream.Recv()
			if err == io.EOF {
				stream.Close()
				close(pipe)
				break
			}
			pipe <- part.Choices[0].Delta.Content
		}
	}()
	return pipe
}

func LogAndGetString(cmd *cobra.Command, c chan string) string {
	var fullRecipe strings.Builder
	var done bool
	for !done {
		select {
		case recipe, ok := <-c:
			if !ok {
				done = true
				break
			}
			cmd.Print(recipe)
			fullRecipe.WriteString(recipe)
		}
	}
	return fullRecipe.String()
}
