package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const openaiApiBase = "https://api.openai.com/v1"

type ChatClient interface {
	RequestCompletion(ctx context.Context, chat CompletionChat) (*Completion, error)
}

type ChatClientImpl struct {
	token      string
	httpClient *http.Client
}

func NewChatClient(token string) *ChatClientImpl {
	return &ChatClientImpl{
		token:      token,
		httpClient: http.DefaultClient,
	}
}

var (
	ErrCompletionFailed = errors.New("completion failed")
	ErrModeration       = errors.New("prompt flagged")
)

// RequestCompletion creates a new chat completion with the given chat configuration.
// User's prompt will be filtered with moderation API. If any of the user messages will be flagged completion will fail with ErrModeration.
func (c *ChatClientImpl) RequestCompletion(ctx context.Context, chat CompletionChat) (*Completion, error) {
	for _, msg := range chat.Messages {
		if msg.Role == "user" {
			result, err := c.moderatePrompt(ctx, msg.Content)
			if err != nil {
				return nil, err
			}

			if result.Results[0].Flagged {
				return nil, ErrModeration
			}
		}
	}

	body, err := json.Marshal(chat)
	if err != nil {
		return nil, err
	}

	req, err := c.request(ctx, http.MethodPost, "/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, ErrCompletionFailed
	}

	var completion Completion
	if err := json.NewDecoder(res.Body).Decode(&completion); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCompletionFailed, err)
	}

	fmt.Println(completion)

	return &completion, nil
}

// moderatePrompt runs OpenAI moderations service on the give prompt.
func (c *ChatClientImpl) moderatePrompt(ctx context.Context, prompt string) (*moderationResult, error) {
	body, err := json.Marshal(moderationRequest{
		Input: prompt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marschal: %w", err)
	}

	req, err := c.request(ctx, http.MethodPost, "/moderations", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("moderation request failed")
	}

	var result moderationResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed decode result: %w", err)
	}

	return &result, nil
}

// request assembles a base request to OpenAI API.
func (c *ChatClientImpl) request(ctx context.Context, method string, endpoint string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, openaiApiBase+endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	return req, nil
}
