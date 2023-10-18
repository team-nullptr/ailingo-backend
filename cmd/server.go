package main

import (
	"ailingo/chat"
	"crypto/tls"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := godotenv.Load(".env"); err != nil {
		l.Error("failed to load configuration from .env file\n")
		os.Exit(1)
	}

	sentenceGenerator := chat.NewSentenceGenerator(http.DefaultClient)
	chatController := chat.NewController(sentenceGenerator)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // Use this to allow specific origin hosts
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// TODO: Use load balancer to
	r.Post("/service/sentence", chatController.GenerateSentence)

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
