package httptransport

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"ailingo/internal/domain"
	"ailingo/pkg/apiutil"
)

type AiController struct {
	logger             *slog.Logger
	chatUseCase        domain.ChatUseCase
	translationUseCase domain.TranslateUseCase
}

func NewAiController(logger *slog.Logger, chatUseCase domain.ChatUseCase, translationUseCase domain.TranslateUseCase) *AiController {
	return &AiController{
		logger:             logger,
		chatUseCase:        chatUseCase,
		translationUseCase: translationUseCase,
	}
}

// GenerateSentence is an endpoint handler for generating a sentence containing submitted word.
func (c *AiController) GenerateSentence(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*30)
	defer cancel()

	var word domain.Definition
	if err := json.NewDecoder(r.Body).Decode(&word); err != nil {
		apiutil.Err(c.logger, w, apiutil.ApiError{
			Status:  http.StatusBadRequest,
			Message: "Invalid sentence generation request payload",
			Cause:   err,
		})
		return
	}

	generatedSentence, err := c.chatUseCase.GenerateSentence(ctx, &word)
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
func (c *AiController) Translate(w http.ResponseWriter, r *http.Request) {
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
