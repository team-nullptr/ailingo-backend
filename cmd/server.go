package main

import (
	"ailingo/internal/chat"
	"ailingo/internal/translation"
	"ailingo/pkg/deepl"
	"ailingo/pkg/openai"
	"crypto/tls"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"os"
)

// TODO: We need a structured logger, look at log/slog package
// TODO: Use load balancer to

func main() {
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := godotenv.Load(".env"); err != nil {
		l.Error("failed to load configuration from .env file\n")
		os.Exit(1)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	openaiClient := openai.NewChatClient(http.DefaultClient, os.Getenv("OPENAI_SECRET"))
	deeplClient := deepl.NewClient(http.DefaultClient, os.Getenv("DEEPL_SECRET"))

	sentenceService := chat.NewSentenceService(openaiClient)
	chatController := chat.NewController(sentenceService)
	chatController.Attach(r, "/gpt")

	translationController := translation.NewController(deeplClient)
	translationController.Attach(r, "/translate")

	srv := http.Server{
		Addr:    os.Getenv("SERVER_ADDR"),
		Handler: r,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	}

	certFile := os.Getenv("TLS_CERT")
	keyFile := os.Getenv("TLS_KEY")

	if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil {
		l.Error("failed to run the server", err)
		os.Exit(1)
	}
}
