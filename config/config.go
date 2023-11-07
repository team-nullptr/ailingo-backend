package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const (
	EnvProd = "PROD"
)

var (
	ErrInvalidValue = fmt.Errorf("invalid value")
)

type Server struct {
	Port               string
	UseTLS             bool
	TLSCert            string
	TLSKey             string
	CorsAllowedOrigins []string
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

	var useTLS bool
	if value := os.Getenv("USE_TLS"); value == "true" {
		useTLS = true
	} else if value == "false" {
		useTLS = false
	} else {
		return nil, fmt.Errorf("%w: invalid value for USE_TLS env variable", ErrInvalidValue)
	}

	return &Config{
		Server: Server{
			Port:               os.Getenv("PORT"),
			UseTLS:             useTLS,
			TLSCert:            os.Getenv("TLS_CERT"),
			TLSKey:             os.Getenv("TLS_KEY"),
			CorsAllowedOrigins: parseOrigins(os.Getenv("CORS_ALLOWED_ORIGINS")),
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

func parseOrigins(origins string) []string {
	return strings.Split(origins, ",")
}
