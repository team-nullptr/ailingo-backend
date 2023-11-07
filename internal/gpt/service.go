package gpt

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"

	"ailingo/internal/domain"
	"ailingo/pkg/openai"
)

const gptModel = "gpt-4-1106-preview"

var (
	// ErrModelDelusions is returned in case of unexpected completion output.
	// As GPT models are not deterministic we cannot assume the output will be always in the form we asked for.
	ErrModelDelusions = errors.New("unexpected completion output")
	// ErrGenerationUnsuccessful is returned in case the model marks their answer as unsuccessful.
	ErrGenerationUnsuccessful = errors.New("generation was not successful")
)

type service struct {
	chatClient openai.ChatClient
}

func NewService(chatClient openai.ChatClient) domain.AiService {
	return &service{
		chatClient: chatClient,
	}
}

// sentenceGeneratorSystem is a prompt for sentence generator persona.
//
//go:embed prompts/sentence_generator.prompt
var sentenceGeneratorSystem string

// SentenceGenerationResult represents GPT model.go response to sentence generation request.
type SentenceGenerationResult struct {
	Success  bool   `json:"success"`
	Sentence string `json:"sentence,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

func (s *service) GenerateSentence(ctx context.Context, req *domain.SentenceGenerationRequest) (string, error) {
	completion, err := s.chatClient.RequestCompletion(ctx, &openai.CompletionChat{
		Model: gptModel,
		Messages: []openai.Message{
			{
				Role:    "system",
				Content: sentenceGeneratorSystem,
			},
			{
				Role: "user",
				Content: fmt.Sprintf(
					"english phrase: %s\npolish meaning: %s",
					req.Phrase,
					req.Meaning,
				),
			},
		},
		MaxTokens: 300,
	})
	if err != nil {
		return "", err
	}

	// TODO: Is it possible to have empty choices array?
	var result SentenceGenerationResult
	if err = json.Unmarshal([]byte(completion.Choices[0].Message.Content), &result); err != nil {
		return "", fmt.Errorf("%w: %w", ErrModelDelusions, err)
	}
	if !result.Success {
		return "", fmt.Errorf("%w: %s", ErrGenerationUnsuccessful, result.Reason)
	}

	return result.Sentence, nil
}

// sentenceGeneratorSystem is a prompt for sentence generator persona.
//
//go:embed prompts/set_generator.prompt
var setGeneratorSystem string

type SetGenerationResult struct {
	Success     bool                           `json:"success"`
	Definitions []*domain.InsertDefinitionData `json:"definitions"`
	Reason      string                         `json:"reason"`
}

func (s *service) GenerateDefinitions(ctx context.Context, req *domain.SetGenerationRequest) ([]*domain.InsertDefinitionData, error) {
	completion, err := s.chatClient.RequestCompletion(ctx, &openai.CompletionChat{
		Model: gptModel,
		Messages: []openai.Message{
			{
				Role:    "system",
				Content: setGeneratorSystem,
			},
			{
				Role:    "user",
				Content: fmt.Sprintf("word_set_title: \"%s\"", req.Name),
			},
		},
		MaxTokens: 1024,
	})
	if err != nil {
		return nil, err
	}

	var result SetGenerationResult
	if err = json.Unmarshal([]byte(completion.Choices[0].Message.Content), &result); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrModelDelusions, err)
	}
	if result.Success {
		return result.Definitions, nil
	}

	return nil, fmt.Errorf("%w: %w", ErrGenerationUnsuccessful, result.Reason)
}
