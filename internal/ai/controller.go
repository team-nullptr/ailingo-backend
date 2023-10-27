package ai

import (
	"ailingo/internal/ai/sentence"
	"ailingo/internal/ai/translate"
	"ailingo/internal/apiutil"

	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"ailingo/internal/models"
)

// Controller exposes handlers for AI related features API.
type Controller struct {
	logger             *slog.Logger
	chatUseCase        sentence.ChatUseCase
	translationUseCase translate.TranslationUseCase
}

func New(logger *slog.Logger, chatUseCase sentence.ChatUseCase, translationUseCase translate.TranslationUseCase) *Controller {
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
		apiutil.Err(c.logger, w, apiutil.ApiError{
			Status:  http.StatusBadRequest,
			Message: "Invalid sentence generation request payload",
			Cause:   err,
		})
		return
	}

	generatedSentence, err := c.chatUseCase.GenerateSentence(ctx, word)
	if err != nil {
		apiutil.Err(c.logger, w, err)
		return
	}

	apiutil.Json(c.logger, w, http.StatusOK, map[string]string{
		"sentence": generatedSentence,
	})
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
		apiutil.Err(c.logger, w, apiutil.ApiError{
			Status:  http.StatusBadRequest,
			Message: "Invalid translation request payload",
			Cause:   err,
		})
		return
	}

	t, err := c.translationUseCase.Translate(ctx, body.Phrase)
	if err != nil {
		apiutil.Err(c.logger, w, err)
		return
	}

	apiutil.Json(c.logger, w, http.StatusOK, map[string]string{
		"definition": t,
	})
}
