package main

import (
	"ailingo/chat"
	"crypto/tls"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("failed to load configuration from .env file\n")
	}

	sentenceGenerator := chat.NewSentenceGenerator(http.DefaultClient)
	chatController := chat.NewController(sentenceGenerator)

	r := chi.NewRouter()
	r.Get("/generation/sentence", chatController.GenerateSentence)

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
		log.Fatalf("failed to start the server: %s\n", err)
	}
}
