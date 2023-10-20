package main

import (
	"ailingo/internal/chat"
	"ailingo/internal/config"
	"ailingo/internal/translation"
	"ailingo/pkg/deepl"
	"ailingo/pkg/openai"
	"crypto/tls"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"log/slog"
	"net/http"
	"os"
)

// TODO: We need a structured logger, look at log/slog package
// TODO: Use load balancer to

func initRouter(cfg *config.Config) *chi.Mux {
	openaiClient := openai.NewChatClient(http.DefaultClient, cfg.OpenaiToken)
	deeplClient := deepl.NewClient(http.DefaultClient, cfg.DeeplToken)
	sentenceService := chat.NewSentenceService(openaiClient)
	chatController := chat.NewController(sentenceService)
	translationController := translation.NewController(deeplClient)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	chatController.Attach(r, "/gpt")
	translationController.Attach(r, "/translate")

	return r
}

func main() {
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		l.Error("failed to load configuration", slog.String("err", err.Error()))
		os.Exit(1)
	}

	r := initRouter(cfg)

	srv := http.Server{
		Addr:    os.Getenv("SERVER_ADDR"),
		Handler: r,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	}

	l.Info(fmt.Sprintf("server starting at %s", cfg.Addr))
	if err := srv.ListenAndServeTLS(cfg.TlsCert, cfg.TlsKey); err != nil {
		l.Error("failed to start the server", slog.String("err", err.Error()))
		os.Exit(1)
	}
}
