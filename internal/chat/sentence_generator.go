package chat

import (
	"ailingo/internal/models"
	"ailingo/pkg/openai"
	"context"
	_ "embed"

	"encoding/json"
	"errors"
	"fmt"
)

//go:embed prompts/sentence_generator_v2.prompt
var __sentenceGeneratorPersona string

var (
	ErrModelDelusions = errors.New("unprocessable output")
)

type SentenceService struct {
	chatClient *openai.ChatClient
}

func NewSentenceService(chatClient *openai.ChatClient) *SentenceService {
	return &SentenceService{
		chatClient: chatClient,
	}
}

type GenerationResult struct {
	Success  bool   `json:"success"`
	Sentence string `json:"sentence,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

// GenerateSentence tries to generate a unique sentence with the provided definition.
func (s *SentenceService) GenerateSentence(ctx context.Context, word models.Word) (*GenerationResult, error) {
	wordPrompt := word.ToChatPrompt()

	completion, err := s.chatClient.RequestCompletion(ctx, openai.CompletionChat{
		Model: "gpt-3.5-turbo",
		Messages: []openai.Message{
			{Role: "system", Content: __sentenceGeneratorPersona},
			{Role: "user", Content: wordPrompt},
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
