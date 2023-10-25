package usecase

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"

	"ailingo/internal/models"
	"ailingo/pkg/openai"
)

var (
	// ErrModelDelusions is returned in case of unexpected completion output.
	// As GPT models are not deterministic we cannot assume the output will be always in the form we asked for.
	ErrModelDelusions = errors.New("unexpected completion output")
)

// sentenceGeneratorSystem is a prompt for sentence generator persona.
//
//go:embed prompts/sentence_generator_v2.prompt
var sentenceGeneratorSystem string

type ChatUseCase interface {
	GenerateSentence(ctx context.Context, word models.Word) (*GenerationResult, error)
}

// ChatUseCaseImpl expose features related with OpenAI's chat completion API.
type ChatUseCaseImpl struct {
	chatClient *openai.ChatClientImpl
}

func NewChat(chatClient *openai.ChatClientImpl) ChatUseCase {
	return &ChatUseCaseImpl{
		chatClient: chatClient,
	}
}

// GenerationResult represents GPT model.go response to sentence generation request.
type GenerationResult struct {
	Success  bool   `json:"success"`
	Sentence string `json:"sentence,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

// GenerateSentence requests a new chat completion with Sentence Generator Persona prompt.
func (uc *ChatUseCaseImpl) GenerateSentence(ctx context.Context, word models.Word) (*GenerationResult, error) {
	completion, err := uc.chatClient.RequestCompletion(ctx, openai.CompletionChat{
		Model: "gpt-3.5-turbo",
		Messages: []openai.Message{
			{Role: "system", Content: sentenceGeneratorSystem},
			{Role: "user", Content: word.ToChatPrompt()},
		},
		MaxTokens: 256,
	})
	if err != nil {
		return nil, err
	}

	// TODO: Is it possible to have empty choices array?
	var genResult GenerationResult
	if err := json.Unmarshal([]byte(completion.Choices[0].Message.Content), &genResult); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrModelDelusions, err)
	}

	return &genResult, nil
}