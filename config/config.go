package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	ENV_PROD = "PROD"
)

// Config stores the app configuration.
type Config struct {
	Env         string
	Port        string
	TlsCert     string
	TlsKey      string
	DSN         string
	ClerkToken  string
	OpenAIToken string
	DeepLToken  string
}

// New loads Config, using .env as the config source, and returns it.
func New(useDotenv bool) (*Config, error) {

	if useDotenv {
		if err := godotenv.Load(".env"); err != nil {
			return nil, fmt.Errorf("failed to load .env: %w", err)
		}
	}

	env := os.Getenv("ENV")
	if env != "PROD" {
		env = "DEV"
	}

	return &Config{
		Env:         env,
		Port:        os.Getenv("PORT"),
		TlsCert:     os.Getenv("TLS_CERT"),
		TlsKey:      os.Getenv("TLS_KEY"),
		DSN:         os.Getenv("DSN"),
		ClerkToken:  os.Getenv("CLERK_TOKEN"),
		OpenAIToken: os.Getenv("OPENAI_TOKEN"),
		DeepLToken:  os.Getenv("DEEPL_TOKEN"),
	}, nil
}
