package deepl

import (
	"bytes"
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

func (c *Client) Translate(text string) (string, error) {
	body, err := json.Marshal(DeeplTranslationRequest{
		Text:       []string{text},
		TargetLang: "PL",
	})
	if err != nil {
		return "", fmt.Errorf("invalid translation request: %w", err)
	}

	req, err := c.request(http.MethodPost, "/translate", bytes.NewReader(body))
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

	var result DeeplTranslationResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode translation response: %w", err)
	}

	return result.Translations[0].Text, nil
}

// Request makes a request to DeepL's free tier api.
func (c *Client) request(method string, endpoint string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, deeplApiBase+endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare the request: %w", err)
	}

	req.Header.Set("Authorization", "DeepL-Auth-Key "+c.token)

	return req, nil
}