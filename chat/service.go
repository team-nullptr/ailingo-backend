package chat

import (
	"ailingo/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const openaiApi = "https://api.openai.com/v1"
const model = "gpt-3.5-turbo"
const sentenceGeneratorPersona = `You will be provided with {{PhraseLang}} phrase and its {{MeaningLang}} meaning. Write a sentence in {{PhraseLang}} using the phrase in its {{MeaningLang}} meaning which can be used for students to understand the phrase better. Do not add any labels, notes or text formatting to your response.
phrase: {{Phrase}}
meaning: {{Meaning}}`

type SentenceGenerator struct {
	httpClient *http.Client
}

func NewSentenceGenerator(httpClient *http.Client) *SentenceGenerator {
	return &SentenceGenerator{
		httpClient: httpClient,
	}
}

// GenerateSentence uses the gpt-3.5-turbo chat api to generate a unique sentence with the provided definition.
func (s *SentenceGenerator) GenerateSentence(definition models.Definition) (string, error) {
	cfg := Chat{
		Model: model,
		Messages: []Message{{
			Role:    "system",
			Content: s.buildPrompt(definition),
		}},
		MaxTokens: 256,
	}

	body, err := json.Marshal(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to marshal chat completion settings: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, openaiApi+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to prepare the generation request: %w", err)
	}

	// TODO: Create global config with used credentials
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_SECRET"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("openai request failed: %w", err)
	}

	var completion Completion
	if err := json.NewDecoder(resp.Body).Decode(&completion); err != nil {
		return "", fmt.Errorf("failed to read the response body: %w", err)
	}

	// TODO: We might want to check the .Choices len before accessing the .Message
	return completion.Choices[0].Message.Content, nil
}

func (s *SentenceGenerator) buildPrompt(definition models.Definition) string {
	return strings.NewReplacer(
		"{{PhraseLang}}", definition.PhraseLang,
		"{{MeaningLang}}", definition.MeaningLang,
		"{{Phrase}}", definition.Phrase,
		"{{Meaning}}", definition.Meaning,
	).Replace(sentenceGeneratorPersona)
}
