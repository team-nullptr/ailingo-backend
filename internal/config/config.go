package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OpenaiToken string
	DeeplToken  string
	Addr        string
	TlsCert     string
	TlsKey      string
	DSN         string
}

func Load() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("failed to load .env: %w", err)
	}

	return &Config{
		OpenaiToken: os.Getenv("OPENAI_TOKEN"),
		DeeplToken:  os.Getenv("DEEPL_TOKEN"),
		Addr:        os.Getenv("SERVER_ADDR"),
		TlsCert:     os.Getenv("TLS_CERT"),
		TlsKey:      os.Getenv("TLS_KEY"),
		DSN:         os.Getenv("DSN"),
	}, nil
}
