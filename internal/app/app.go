package app

import (
	"context"
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
	"ailingo/internal/domain"
	"ailingo/internal/gpt"
	"ailingo/internal/mysql"
	"ailingo/internal/usecase"
	"ailingo/internal/webhook"
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
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Database
	// TODO: Implement connection retries
	db, err := connectToDatabase(cfg.Database.DSN)
	if err != nil {
		l.Error(fmt.Sprintf("app - Run - connectToDatabase: %s", err))
		os.Exit(1)
	}

	// Clerk
	clerkClient, err := clerk.NewClient(cfg.Services.ClerkToken)
	if err != nil {
		l.Error(fmt.Sprintf("app - Run - clerk.NewClient: %s", err))
		os.Exit(1)
	}

	withClaims := auth.WithClaims(l, clerkClient)
	userService := auth.NewUserService(l, clerkClient)

	// Repos
	mysqlDataStore := mysql.NewDataStore(db)

	// TODO: This is stupid
	if err := mysqlDataStore.Atomic(context.Background(), func(ds domain.DataStore) error {
		users, err := clerkClient.Users().ListAll(clerk.ListAllUsersParams{})
		if err != nil {
			return fmt.Errorf("failed to list all users: %w", err)
		}

		if err := ds.GetUserRepo().SyncUsers(context.Background(), users); err != nil {
			return fmt.Errorf("failed to sync users: %w", err)
		}

		return nil
	}); err != nil {
		l.Error(fmt.Sprintf("app - Run - user sync: %s", err))
		os.Exit(1)
	} else {
		l.Info("Users have been successfully synced")
	}

	sentenceRepo := gpt.NewSentenceDevRepo()
	if cfg.Server.Env == config.EnvProd {
		sentenceRepo = gpt.NewSentenceRepo(openai.NewChatClient(cfg.Services.OpenAIToken))
	}

	// Validator
	validate := validator.New(validator.WithRequiredStructEnabled())

	// Use cases
	translationUseCase := usecase.NewTranslateDevUseCase()
	if cfg.Server.Env == config.EnvProd {
		translationUseCase = usecase.NewTranslateUseCase(deepl.NewClient(cfg.Services.DeepLToken), validate)
	}
	chatUseCase := usecase.NewChatUseCase(sentenceRepo, validate)

	studySetUseCase := usecase.NewStudySetUseCase(mysqlDataStore, userService, validate)
	definitionUseCase := usecase.NewDefinitionUseCase(mysqlDataStore, validate)
	profileUseCase := usecase.NewProfileUseCase(mysqlDataStore, userService)
	userUseCase := usecase.NewUserUseCase(mysqlDataStore)
	studySessionUseCase := usecase.NewStudySessionUseCase(mysqlDataStore)

	// Controllers
	ai := controller.NewAiController(
		l,
		chatUseCase,
		translationUseCase,
	)

	studySet := controller.NewStudySetController(
		l,
		userService,
		studySetUseCase,
		definitionUseCase,
	)

	me := controller.NewMeController(l, profileUseCase, studySessionUseCase, userService)

	clerkWebhook, err := webhook.NewClerkWebhook(l, cfg, userUseCase)
	if err != nil {
		l.Error(fmt.Sprintf("app - Run - webhook.NewClerkWebhook: %s", err))
		os.Exit(1)
	}

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
		AllowedOrigins:   cfg.Server.CorsAllowedOrigins,
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

	r.Route("/study-sets", studySet.Router(withClaims))
	r.With(withClaims).Route("/ai", ai.Router)
	r.With(withClaims).Route("/me", me.Router)

	// Webhooks
	r.Post("/clerk/webhook", clerkWebhook.Webhook)

	// Server
	server := httpserver.New(
		httpserver.WithAddr(fmt.Sprintf(":%s", cfg.Server.Port)),
		httpserver.WithReadTimeout(5*time.Second),
		httpserver.WithWriteTimeout(10*time.Second),
		httpserver.WithHandler(r),
	)

	if cfg.Server.UseTLS {
		server.StartTLS(cfg.Server.TLSCert, cfg.Server.TLSKey)
	} else {
		server.Start()
	}

	// Interrupt signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info(fmt.Sprintf("app - Run - signal: %s", s.String()))
	case err := <-server.Notify():
		l.Error(fmt.Sprintf("app - Run - httpServer.Notify: %s", err))
	}

	// Graceful shutdown
	if err := server.Shutdown(); err != nil {
		l.Error(fmt.Sprintf("app - Run - httpServer.Shutdown: %s", err))
	}
}
