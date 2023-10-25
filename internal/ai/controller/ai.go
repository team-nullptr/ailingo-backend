package controller

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	usecase2 "ailingo/internal/ai/usecase"
	"ailingo/internal/models"
	"ailingo/pkg/apiutil"
)

// Controller exposes handlers for AI related features API.
type Controller struct {
	logger             *slog.Logger
	chatUseCase        usecase2.ChatUseCase
	translationUseCase usecase2.TranslationUseCase
}

func New(logger *slog.Logger, chatUseCase usecase2.ChatUseCase, translationUseCase usecase2.TranslationUseCase) *Controller {
	return &Controller{
		logger:             logger,
		chatUseCase:        chatUseCase,
		translationUseCase: translationUseCase,
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

	result, err := c.chatUseCase.GenerateSentence(ctx, word)
	if err != nil {
		apiutil.Err(c.logger, w, http.StatusInternalServerError, err)
		return
	}

	apiutil.Json(c.logger, w, http.StatusOK, result)
}

type translatePayload struct {
	Phrase string `json:"phrase"`
}

// Translate is an endpoint handler for translating words using DeepL.
func (c *Controller) Translate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*30)
	defer cancel()

	var body translatePayload
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apiutil.Err(c.logger, w, http.StatusBadRequest, errors.New("unprocessable request body"))
		return
	}

	t, err := c.translationUseCase.Translate(ctx, body.Phrase)
	if err != nil {
		apiutil.Err(c.logger, w, http.StatusInternalServerError, err)
		return
	}

	apiutil.Json(c.logger, w, http.StatusOK, map[string]string{
		"definition": t,
	})
}