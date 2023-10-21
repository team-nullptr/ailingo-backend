package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/go-sql-driver/mysql"

	"ailingo/internal/chat"
	"ailingo/internal/config"
	"ailingo/internal/translation"
	"ailingo/pkg/deepl"
	"ailingo/pkg/openai"
)

// TODO: We need a structured logger, look at log/slog package.
// TODO: Should we consider implementing load balancer?

func initRouter(cfg *config.Config) (*chi.Mux, error) {
	_, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	translator := translation.NewDevTranslator()
	chatService := chat.NewDevService()

	if cfg.Prod {
		translator = deepl.NewClient(http.DefaultClient, cfg.DeepLToken)
		openaiClient := openai.NewChatClient(http.DefaultClient, cfg.OpenAIToken)
		chatService = chat.NewService(openaiClient)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	chat.NewController(chatService).Attach(r, "/gpt")
	translation.NewController(translator).Attach(r, "/translate")

	return r, nil
}

func main() {
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		l.Error("failed to load configuration", slog.String("err", err.Error()))
		os.Exit(1)
	}

	l.Info(fmt.Sprintf("production %t", cfg.Prod))

	r, err := initRouter(cfg)
	if err != nil {
		l.Error("failed to initialize server router", slog.String("err", err.Error()))
		os.Exit(1)
	}

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
