package deepl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const deeplApiBase = "https://api-free.deepl.com/v2"

type Client struct {
	token      string
	httpClient *http.Client
}

func NewClient(httpClient *http.Client, token string) *Client {
	return &Client{
		token:      token,
		httpClient: httpClient,
	}
}

// Translate translates given text into english using DeepL API.
func (c *Client) Translate(ctx context.Context, text string) (string, error) {
	body, err := json.Marshal(TranslationRequest{
		Text:       []string{text},
		TargetLang: "PL",
	})
	if err != nil {
		return "", fmt.Errorf("invalid translation request: %w", err)
	}

	req, err := c.request(ctx, http.MethodPost, "/translate", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("translation error: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("translation error failed %d", res.StatusCode)
	}

	var result TranslationResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode translation response: %w", err)
	}

	return result.Translations[0].Text, nil
}

// request assembles a base request to DeepL's free tier api.
func (c *Client) request(ctx context.Context, method string, endpoint string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, deeplApiBase+endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare the request: %w", err)
	}

	req.Header.Set("Authorization", "DeepL-Auth-Key "+c.token)

	return req, nil
}
