package app

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/httprate"
	"github.com/go-playground/validator/v10"

	"ailingo/config"
	"ailingo/internal/controller"
	"ailingo/internal/gpt"
	"ailingo/internal/mysql"
	"ailingo/internal/usecase"
	"ailingo/pkg/auth"
	"ailingo/pkg/deepl"
	"ailingo/pkg/httpserver"
	"ailingo/pkg/openai"
)

// connectToDatabase creates a connection to the database via the given dataSourceName.
func connectToDatabase(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Run starts the application's.
func Run(cfg *config.Config) {
	// Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Database
	// TODO: Implement connection retries
	db, err := connectToDatabase(cfg.DSN)
	if err != nil {
		logger.Error(fmt.Sprintf("app - Run - connectToDatabase: %s", err))
		os.Exit(1)
	}

	// Clerk
	clerkClient, err := clerk.NewClient(cfg.ClerkToken)
	if err != nil {
		logger.Error(fmt.Sprintf("app - Run - clerk.NewClient: %s", err))
		os.Exit(1)
	}

	withClaims := auth.WithClaims(logger, clerkClient)
	userService := auth.NewUserService(logger, clerkClient)

	// Repos
	sentenceRepo := gpt.NewSentenceDevRepo()
	if cfg.Env == config.ENV_PROD {
		sentenceRepo = gpt.NewSentenceRepo(openai.NewChatClient(cfg.OpenAIToken))
	}

	definitionRepo, err := mysql.NewDefinitionRepo(db)
	if err != nil {
		logger.Error(fmt.Sprintf("app - Run - mysql.NewDefinitionRepo: %s", err))
		os.Exit(1)
	}

	studySetRepo, err := mysql.NewStudySetRepo(db)
	if err != nil {
		logger.Error(fmt.Sprintf("app - Run - mysql.NewStudySetRepo: %s", err))
		os.Exit(1)
	}

	// Validator
	validate := validator.New(validator.WithRequiredStructEnabled())

	// Use cases
	translationUseCase := usecase.NewTranslateDevUseCase()
	if cfg.Env == config.ENV_PROD {
		translationUseCase = usecase.NewTranslateUseCase(deepl.NewClient(cfg.DeepLToken), validate)
	}
	chatUseCase := usecase.NewChatUseCase(sentenceRepo, validate)
	studySetUseCase := usecase.NewStudySetUseCase(studySetRepo, definitionRepo, validate)
	definitionUseCase := usecase.NewDefinitionUseCase(definitionRepo, studySetRepo)

	// Controllers
	ai := controller.NewAiController(
		logger,
		chatUseCase,
		translationUseCase,
	)

	studySet := controller.NewStudySetController(
		logger,
		userService,
		studySetUseCase,
		definitionUseCase,
	)

	// Router
	reqLogger := httplog.RequestLogger(httplog.NewLogger("api", httplog.Options{
		LogLevel:      slog.LevelDebug,
		JSON:          true,
		TimeFieldName: "time",
	}))

	limiter := httprate.Limit(
		10,
		10*time.Second,
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
	)

	corsOpts := cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	r := chi.NewRouter()
	r.Use(
		limiter,
		reqLogger,
		corsOpts,
	)
	r.With(withClaims).Route("/ai", ai.Router)
	r.Route("/study-sets", studySet.Router(withClaims))

	// Server
	server := httpserver.New(
		httpserver.WithHandler(r),
		httpserver.WithAddr(fmt.Sprintf(":%s", cfg.Port)),
	)

	server.Start(cfg.TlsCert, cfg.TlsKey)

	// Interrupt signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info(fmt.Sprintf("app - Run - signal: %s", s.String()))
	case err := <-server.Notify():
		logger.Error(fmt.Sprintf("app - Run - httpServer.Notify: %s", err))
	}

	// Graceful shutdown
	if err := server.Shutdown(); err != nil {
		logger.Error(fmt.Sprintf("app - Run - httpServer.Shutdown: %s", err))
	}
}
