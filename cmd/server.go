package main

import (
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

	"ailingo/internal/chat"
	"ailingo/internal/config"
	"ailingo/internal/studyset"
	"ailingo/internal/translation"
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

	var translator translation.Translator
	var chatService chat.Service

	if cfg.Env == config.ENV_PROD {
		translator = deepl.NewClient(cfg.DeepLToken)
		chatService = chat.NewService(openai.NewChatClient(cfg.OpenAIToken))
	} else {
		translator = translation.NewDevTranslator()
		chatService = chat.NewDevService()
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	studySetService := studyset.NewService(studyset.NewRepo(db), validate)

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
		chatController := chat.NewController(logger, chatService)
		translationController := translation.NewController(logger, translator)

		r.Post("/sentence", chatController.GenerateSentence)
		r.Post("/translate", translationController.Translate)
	})

	r.Route("/study-sets", func(r chi.Router) {
		studysetController := studyset.NewController(logger, studySetService)

		r.Use(withClaims)
		r.Post("/", studysetController.Create)
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
