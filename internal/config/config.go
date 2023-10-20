package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	OpenaiToken string
	DeeplToken  string
	Addr        string
	TlsCert     string
	TlsKey      string
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
	}, nil
}
