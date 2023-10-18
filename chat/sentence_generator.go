package chat

import (
	_ "embed"

	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"ailingo/models"
)

const gptModel = "gpt-3.5-turbo"

//go:embed prompts/sentence_generator_v2.prompt
var __sentenceGeneratorPersona string

var (
	ErrModeration     = errors.New("prompt has been flagged by the moderation service")
	ErrModelDelusions = errors.New("unprocessable output")
)

type SentenceGenerator struct {
	httpClient *http.Client
	// TODO: We need a structured logger, look at log/slog package
}

// NewSentenceGenerator creates a new sentence generator service.
func NewSentenceGenerator(httpClient *http.Client) *SentenceGenerator {
	return &SentenceGenerator{
		httpClient: httpClient,
	}
}

type GenerationResult struct {
	Success  bool   `json:"success"`
	Sentence string `json:"sentence,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

// GenerateSentence tries to generate a unique sentence with the provided definition.
func (s *SentenceGenerator) GenerateSentence(word models.Word) (*GenerationResult, error) {
	wordPrompt := word.ToChatPrompt()

	moderation, err := s.runModeration(wordPrompt)
	if err != nil {
		return nil, fmt.Errorf("moderation apiutil failed: %w", err)
	}
	if len(moderation.Results) > 0 && moderation.Results[0].Flagged {
		return nil, ErrModeration
	}

	body, err := json.Marshal(CompletionChat{
		Model: gptModel,
		Messages: []Message{
			{Role: "system", Content: __sentenceGeneratorPersona},
			{Role: "user", Content: wordPrompt},
		},
		MaxTokens: 256,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chat completion settings: %w", err)
	}

	req, err := newApiRequest(http.MethodPost, "/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("completion request failed: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("completion failed")
	}

	var completion Completion
	if err := json.NewDecoder(res.Body).Decode(&completion); err != nil {
		return nil, fmt.Errorf("failed to decode completion response: %w", err)
	}

	// TODO: Is it possible to have empty choices array?
	var genResult GenerationResult
	if err := json.Unmarshal([]byte(completion.Choices[0].Message.Content), &genResult); err != nil {
		return nil, fmt.Errorf("%w: failed to unmarshal generation result: %w", ErrModelDelusions, err)
	}

	return &genResult, nil
}

type moderationRequest struct {
	Input string `json:"input"`
}

type moderationResult struct {
	Results []struct {
		Flagged        bool               `json:"flagged"`
		Categories     map[string]bool    `json:"categories"`
		CategoryScores map[string]float32 `json:"category_scores"`
	} `json:"results"`
}

func (s *SentenceGenerator) runModeration(prompt string) (*moderationResult, error) {
	body, err := json.Marshal(moderationRequest{
		Input: prompt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marschal: %w", err)
	}

	req, err := newApiRequest(http.MethodPost, "/moderations", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("moderation request failed")
	}

	var result moderationResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode moderation result: %w", err)
	}

	fmt.Printf("%+v", result)
	return &result, nil
}
