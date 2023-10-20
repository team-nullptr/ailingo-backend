package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const openaiApiBase = "https://api.openai.com/v1"

type ChatClient struct {
	token      string
	httpClient *http.Client
}

func NewChatClient(httpClient *http.Client, token string) *ChatClient {
	return &ChatClient{
		token:      token,
		httpClient: httpClient,
	}
}

var (
	ErrCompletionFailed = errors.New("completion failed")
	ErrModeration       = errors.New("prompt flagged")
)

func (cc *ChatClient) RequestCompletion(chat CompletionChat) (*Completion, error) {
	for _, msg := range chat.Messages {
		if msg.Role == "user" {
			result, err := cc.moderatePrompt(msg.Content)
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

	req, err := cc.request(http.MethodPost, "/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := cc.httpClient.Do(req)
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

	return &completion, nil
}

// moderatePrompt runs OpenAI's moderation service on the give prompt.
func (cc *ChatClient) moderatePrompt(prompt string) (*moderationResult, error) {
	body, err := json.Marshal(moderationRequest{
		Input: prompt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marschal: %w", err)
	}

	req, err := cc.request(http.MethodPost, "/moderations", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	res, err := cc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("moderation request failed")
	}

	var result moderationResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed decode result: %w", err)
	}

	return &result, nil
}

func (cc *ChatClient) request(method string, endpoint string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, openaiApiBase+endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+cc.token)

	return req, nil
}
