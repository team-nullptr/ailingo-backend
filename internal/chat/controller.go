package chat

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"ailingo/internal/models"
	"ailingo/pkg/apiutil"
)

// Controller exposes handlers for GPT API.
type Controller struct {
	chatSvc Service
	logger  *slog.Logger
}

func NewController(logger *slog.Logger, chatSvc Service) *Controller {
	return &Controller{
		chatSvc: chatSvc,
		logger:  logger,
	}
}

// Attach attaches controller to the given mux.
func (c *Controller) Attach(m *chi.Mux, path string) {
	m.Route(path, func(r chi.Router) {
		r.Post("/sentence", c.GenerateSentence)
	})
}

// GenerateSentence is an endpoint handler for generating a sentence containing submitted word.
func (c *Controller) GenerateSentence(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*30)
	defer cancel()

	var word models.Word
	if err := json.NewDecoder(r.Body).Decode(&word); err != nil {
		apiutil.Err(c.logger, w, http.StatusBadRequest, err)
		return
	}

	result, err := c.chatSvc.GenerateSentence(ctx, word)
	if err != nil {
		apiutil.Err(c.logger, w, http.StatusInternalServerError, err)
		return
	}

	apiutil.Json(c.logger, w, http.StatusOK, result)
}
