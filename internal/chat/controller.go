package chat

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"ailingo/internal/models"
	"ailingo/pkg/apiutil"
)

// Controller exposes handlers for GPT API.
type Controller struct {
	chatService Service
	logger      *slog.Logger
}

func NewController(logger *slog.Logger, chatService Service) *Controller {
	return &Controller{
		chatService: chatService,
		logger:      logger,
	}
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

	result, err := c.chatService.GenerateSentence(ctx, word)
	if err != nil {
		apiutil.Err(c.logger, w, http.StatusInternalServerError, err)
		return
	}

	apiutil.Json(c.logger, w, http.StatusOK, result)
}
