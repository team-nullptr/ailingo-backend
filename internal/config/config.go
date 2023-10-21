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
	Port        string
	TlsCert     string
	TlsKey      string
	DSN         string
	Env         string
}

// Load loads Config, using .env as the config source, and returns it.
func Load() (*Config, error) {
	dotenv := flag.Bool("dotenv", false, "configure with .env")

	if *dotenv {
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
		TlsCert:     "./certs/" + os.Getenv("TLS_CERT"),
		TlsKey:      "./certs/" + os.Getenv("TLS_KEY"),
		DSN:         os.Getenv("DSN"),
		OpenAIToken: os.Getenv("OPENAI_TOKEN"),
		DeepLToken:  os.Getenv("DEEPL_TOKEN"),
	}, nil
}
