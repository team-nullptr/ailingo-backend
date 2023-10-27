package sentence

import (
	"ailingo/internal/models"
	"ailingo/pkg/openai"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
)

// sentenceGeneratorSystem is a prompt for sentence generator persona.
//
//go:embed prompts/sentence_generator_v2.prompt
var sentenceGeneratorSystem string

var (
	// ErrModelDelusions is returned in case of unexpected completion output.
	// As GPT models are not deterministic we cannot assume the output will be always in the form we asked for.
	ErrModelDelusions = errors.New("unexpected completion output")

	// ErrGenerationUnsuccessful is returned in case the model marks their answer as unsuccessful.
	ErrGenerationUnsuccessful = errors.New("generation was not successful")
)

type Repo interface {
	GenerateSentence(ctx context.Context, word models.Word) (string, error)
}

type DefaultRepo struct {
	chatClient openai.ChatClient
}

func NewRepo(chatClient openai.ChatClient) Repo {
	return &DefaultRepo{
		chatClient: chatClient,
	}
}

// GenerationResult represents GPT model.go response to sentence generation request.
type GenerationResult struct {
	Success  bool   `json:"success"`
	Sentence string `json:"sentence,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

// GenerateSentence requests a new chat completion with the Sentence Generator Persona prompt.
func (r DefaultRepo) GenerateSentence(ctx context.Context, word models.Word) (string, error) {
	completion, err := r.chatClient.RequestCompletion(ctx, openai.CompletionChat{
		Model: "gpt-3.5-turbo",
		Messages: []openai.Message{
			{Role: "system", Content: sentenceGeneratorSystem},
			{Role: "user", Content: word.ToChatPrompt()},
		},
		MaxTokens: 256,
	})
	if err != nil {
		return "", err
	}
	completionContent := completion.Choices[0].Message.Content
	// TODO: Is it possible to have empty choices array?
	var genResult GenerationResult
	err = json.Unmarshal([]byte(completionContent), &genResult)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrModelDelusions, err)
	}
	if genResult.Success {
		return genResult.Sentence, nil
	}
	return "", fmt.Errorf("%w: %s", ErrGenerationUnsuccessful, genResult.Reason)
}
