package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config stores the app configuration.
type Config struct {
	OpenAIToken string
	DeepLToken  string
	Addr        string
	TlsCert     string
	TlsKey      string
	DSN         string
	Prod        bool
}

// Load loads Config, using .env as the config source, and returns it.
func Load() (*Config, error) {
	prodFlag := flag.Bool("prod", false, "run server in production configuration")
	flag.Parse()

	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("failed to load .env: %w", err)
	}

	return &Config{
		OpenAIToken: os.Getenv("OPENAI_TOKEN"),
		DeepLToken:  os.Getenv("DEEPL_TOKEN"),
		Addr:        os.Getenv("SERVER_ADDR"),
		TlsCert:     os.Getenv("TLS_CERT"),
		TlsKey:      os.Getenv("TLS_KEY"),
		DSN:         os.Getenv("DSN"),
		Prod:        *prodFlag,
	}, nil
}
