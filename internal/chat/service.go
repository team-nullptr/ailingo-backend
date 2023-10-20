package chat

import (
	"context"
	_ "embed"

	"ailingo/internal/models"
	"ailingo/pkg/openai"

	"encoding/json"
	"errors"
	"fmt"
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

// Service expose features related with OpenAI's chat completion API.
type Service struct {
	chatClient *openai.ChatClient
}

func NewService(chatClient *openai.ChatClient) *Service {
	return &Service{
		chatClient: chatClient,
	}
}

// GenerationResult represents GPT model response to sentence generation request.
type GenerationResult struct {
	Success  bool   `json:"success"`
	Sentence string `json:"sentence,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

// GenerateSentence requests a new chat completion with Sentence Generator Persona prompt.
func (s *Service) GenerateSentence(ctx context.Context, word models.Word) (*GenerationResult, error) {
	completion, err := s.chatClient.RequestCompletion(ctx, openai.CompletionChat{
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
