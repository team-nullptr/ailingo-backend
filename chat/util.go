package chat

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

const apiBase = "https://api.openai.com/v1"

func newApiRequest(method string, endpoint string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, apiBase+endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare the request: %w", err)
	}

	// TODO: Create global config with used credentials
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_SECRET"))
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
