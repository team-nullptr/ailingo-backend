package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/httprate"

	"ailingo/internal/auth"
	"ailingo/internal/gpt"
	"ailingo/internal/httptransport"
	"ailingo/internal/mysql"
	"ailingo/internal/usecase"
	"ailingo/pkg/deepl"
	"ailingo/pkg/openai"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"

	"ailingo/internal/config"
)

type ConnectionSettings struct {
	Attempts int
	Sleep    time.Duration
}

// connectToDb establishes a new connection with the database.
func connectToDb(logger *slog.Logger, cfg *config.Config, settings *ConnectionSettings) (*sql.DB, error) {
	var (
		db  *sql.DB
		err error
	)

	for i := 0; i < settings.Attempts; i++ {
		db, err = sql.Open("mysql", cfg.DSN)
		if err == nil {
			if err = db.Ping(); err == nil {
				return db, nil
			}
		}
		if i > 0 {
			logger.Info("Db connection failed", slog.Int("attempt", i), slog.Float64("retryIn", settings.Sleep.Seconds()))
			time.Sleep(settings.Sleep)
			settings.Sleep += time.Second * 2
		}
	}

	return nil, fmt.Errorf("failed to connect with the database after %d attempts: %w", settings.Attempts, err)
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load configuration", slog.String("err", err.Error()))
		os.Exit(1)
	}
	logger.Info("starting in " + cfg.Env + " environment")

	db, err := connectToDb(logger, cfg, &ConnectionSettings{
		Attempts: 10,
		Sleep:    time.Second * 2,
	})
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

	translationUseCase := usecase.NewTranslateDevUseCase()
	sentenceRepo := gpt.NewSentenceDevRepo()
	if cfg.Env == config.ENV_PROD {
		translationUseCase = usecase.NewTranslateUseCase(deepl.NewClient(cfg.DeepLToken))
		sentenceRepo = gpt.NewSentenceRepo(openai.NewChatClient(cfg.OpenAIToken))
	}

	studySetRepo, err := mysql.NewStudySetRepo(db)
	if err != nil {
		logger.Error("failed to initialize study set repo", slog.String("err", err.Error()))
		os.Exit(1)
	}

	chatUseCase := usecase.NewChatUseCase(sentenceRepo)
	studySetUseCase := usecase.NewStudySetUseCase(studySetRepo, validate)

	userService := auth.NewUserService(logger, clerkClient)

	withClaims := auth.WithClaims(logger, clerkClient)
	r := chi.NewRouter()

	r.Use(httprate.Limit(
		10,
		10*time.Second,
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
	))

	r.Use(httplog.RequestLogger(
		httplog.NewLogger("api", httplog.Options{
			LogLevel:      slog.LevelDebug,
			JSON:          true,
			TimeFieldName: "time",
		}),
	))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.With(withClaims).Route("/ai", func(r chi.Router) {
		c := httptransport.NewAiController(logger, chatUseCase, translationUseCase)
		r.Post("/sentence", c.GenerateSentence)
		r.Post("/translate", c.Translate)
	})

	r.Route("/study-sets", func(r chi.Router) {
		c := httptransport.NewStudySetController(logger, userService, studySetUseCase)
		r.Get("/", c.GetAllSummary)
		r.Get("/{studySetID}", c.GetById)
		r.With(withClaims).Post("/", c.Create)
		r.With(withClaims).Put("/{studySetID}", c.Update)
		r.With(withClaims).Delete("/{studySetID}", c.Delete)
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
