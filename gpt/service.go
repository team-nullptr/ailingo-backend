package gpt

import (
	"ailingo/chat"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

// TODO: Maybe we can use of openai moderate endpoint to make sure nobody tries to submit offensive content that could break our model?

const openaiApi = "https://api.openai.com/v1"

var (
	ErrInternal = errors.New("internal error")
	ErrOpenai   = errors.New("request to the openai api failed")
)

type Definition struct {
	Phrase  string
	Meaning string
}

type ExampleSentenceService struct {
	httpClient *http.Client
}

func NewExampleSentenceService(httpClient *http.Client) *ExampleSentenceService {
	return &ExampleSentenceService{
		httpClient: httpClient,
	}
}

var sentenceGeneratorPersona = "You will be provided with english phrase and its polish meaning. " +
	"Write an english sentence using the phrase in its polish meaning which can be used for students to understand the phrase better. " +
	"Do not add any labels, notes or text formatting to your response."

// GenerateSentence uses the gpt api to generate a unique example message for the given definition.
func (s *ExampleSentenceService) GenerateSentence(definition Definition) (string, error) {
	// Prepare prompt

	prompt := sentenceGeneratorPersona +
		"\nphrase: " + definition.Phrase +
		"\nmeaning: " + definition.Meaning

	body, err := json.Marshal(
		chat.CompletionSettings{
			Model:     "gpt-3.5-turbo",
			Messages:  []chat.Message{{Role: "system", Content: prompt}},
			MaxTokens: 256,
		},
	)
	if err != nil {
		return "", fmt.Errorf("%w: failed to marshal chat completion settings: %w", ErrInternal, err)
	}

	// Api call

	req, err := http.NewRequest(http.MethodPost, openaiApi+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("%w: failed to prepare the generation request: %w", ErrInternal, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_SECRET"))

	fmt.Println(req.Header)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrOpenai, err)
	}

	msg, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%w: failed to read the response body: %w", ErrInternal, err)
	}

	fmt.Println(resp.Status, string(msg))

	//var completion chat.Completion
	//if err := json.NewDecoder(resp.Body).Decode(&completion); err != nil {
	//	return "", fmt.Errorf("%w: failed to read the response body: %w", ErrInternal, err)
	//}
	//
	//fmt.Printf("%+v\n", completion)
	return "", nil
}
