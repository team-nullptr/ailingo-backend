package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	EnvProd = "PROD"
)

type Server struct {
	Env     string
	Port    string
	TlsCert string
	TlsKey  string
}

type Database struct {
	DSN string
}

type Services struct {
	DeepLToken         string
	OpenAIToken        string
	ClerkToken         string
	ClerkWebhookSecret string
}

// Config stores the app configuration.
type Config struct {
	Server   Server
	Database Database
	Services Services
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
		Server: Server{
			Env:     env,
			Port:    os.Getenv("PORT"),
			TlsCert: os.Getenv("TLS_CERT"),
			TlsKey:  os.Getenv("TLS_KEY"),
		},
		Database: Database{
			DSN: os.Getenv("DSN"),
		},
		Services: Services{
			ClerkToken:         os.Getenv("CLERK_TOKEN"),
			ClerkWebhookSecret: os.Getenv("CLERK_WEBHOOK_SECRET"),
			OpenAIToken:        os.Getenv("OPENAI_TOKEN"),
			DeepLToken:         os.Getenv("DEEPL_TOKEN"),
		},
	}, nil
}
