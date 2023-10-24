package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"

	"ailingo/internal/chat"
	"ailingo/internal/config"
	"ailingo/internal/studyset"
	"ailingo/internal/translation"
	"ailingo/pkg/deepl"
	"ailingo/pkg/openai"
)

// connectToDatabase establishes a new connection with the database.
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

// newRouter assembles the API router.
func newRouter(l *slog.Logger, cfg *config.Config, db *sql.DB) (*chi.Mux, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	translator := translation.NewDevTranslator()
	chatService := chat.NewDevService()

	studySetRepo := studyset.NewRepo(db)
	studySetService := studyset.NewService(studySetRepo, validate)

	// If the app is running in PROD mode substitute service stubs with production ones.
	if cfg.Env == "PROD" {
		translator = deepl.NewClient(http.DefaultClient, cfg.DeepLToken)
		openaiClient := openai.NewChatClient(http.DefaultClient, cfg.OpenAIToken)
		chatService = chat.NewService(openaiClient)
	}

	r := chi.NewRouter()

	// TODO: Implement load balancer?

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

	chat.NewController(l, chatService).Attach(r, "/gpt")
	translation.NewController(l, translator).Attach(r, "/translate")
	studyset.NewController(l, studySetService).Attach(r, "/study-sets")

	return r, nil
}

func main() {
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Config
	cfg, err := config.Load()
	if err != nil {
		l.Error("failed to load configuration", slog.String("err", err.Error()))
		os.Exit(1)
	}

	l.Info("starting in " + cfg.Env + " environment")

	// Database connection
	db, err := connectToDb(cfg)
	if err != nil {
		l.Error("failed to establish db connection", slog.String("err", err.Error()))
		os.Exit(1)
	}

	l.Info("connected to the database")

	// Router
	r, err := newRouter(l, cfg, db)
	if err != nil {
		l.Error("failed to create server router", slog.String("err", err.Error()))
		os.Exit(1)
	}

	srv := http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	}

	l.Info(fmt.Sprintf("server starting at port %s", cfg.Port))
	if err := srv.ListenAndServeTLS(cfg.TlsCert, cfg.TlsKey); err != nil {
		l.Error("failed to start the server", slog.String("err", err.Error()))
		os.Exit(1)
	}
}
