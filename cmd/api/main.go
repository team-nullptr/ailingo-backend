package main

import (
	"ailingo/internal/ai"
	"ailingo/internal/ai/sentence"
	"ailingo/internal/ai/translate"
	"crypto/tls"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"

	"ailingo/internal/config"
	"ailingo/internal/studyset"
	"ailingo/pkg/auth"
	"ailingo/pkg/deepl"
	"ailingo/pkg/openai"
)

// connectToDb establishes a new connection with the database.
func connectToDb(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping the db: %w", err)
	}

	return db, nil
}

func main() {
	// loggers, connections, configs

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load configuration", slog.String("err", err.Error()))
		os.Exit(1)
	}
	logger.Info("starting in " + cfg.Env + " environment")

	db, err := connectToDb(cfg)
	if err != nil {
		logger.Error("failed to establish db connection", slog.String("err", err.Error()))
		os.Exit(1)
	}
	logger.Info("connected to the database")

	clerkClient, err := clerk.NewClient(cfg.ClerkToken)
	if err != nil {
		logger.Error("failed to create clerk client")
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	// repos, use-cases

	translationRepo := translate.NewDevRepo()
	sentenceRepo := sentence.NewDevRepo()
	if cfg.Env == config.ENV_PROD {
		translationRepo = translate.NewRepo(deepl.NewClient(cfg.DeepLToken))
		sentenceRepo = sentence.NewRepo(openai.NewChatClient(cfg.OpenAIToken))
	}
	translationUseCase := translate.NewTranslationUseCase(translationRepo)
	chatUseCase := sentence.NewChatUseCase(sentenceRepo)

	studySetRepo, err := studyset.NewRepo(db)
	if err != nil {
		logger.Error("failed to initialize study set repo", slog.String("err", err.Error()))
		os.Exit(1)
	}
	studySetUseCase := studyset.NewUseCase(studySetRepo, validate)

	// app router

	withClaims := auth.WithClaims(logger, clerkClient)
	r := chi.NewRouter()

	r.Use(httplog.RequestLogger(
		httplog.NewLogger("api", httplog.Options{
			LogLevel: slog.LevelDebug,
			Concise:  true,
			JSON:     true,
		}),
	))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.Route("/ai", func(r chi.Router) {
		c := ai.New(logger, chatUseCase, translationUseCase)
		r.Post("/sentence", c.GenerateSentence)
		r.Post("/translate", c.Translate)
	})

	r.Route("/study-sets", func(r chi.Router) {
		c := studyset.New(logger, studySetUseCase)
		r.Use(withClaims)
		r.Post("/", c.Create)
	})

	server := http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	}

	logger.Info(fmt.Sprintf("server starting at port %s", cfg.Port))

	if err := server.ListenAndServeTLS(cfg.TlsCert, cfg.TlsKey); err != nil {
		logger.Error("failed to start the server", slog.String("err", err.Error()))
		os.Exit(1)
	}
}
