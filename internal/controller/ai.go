package controller

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"ailingo/internal/domain"
	"ailingo/internal/usecase"
	"ailingo/pkg/apiutil"
)

type AiController struct {
	l                  *slog.Logger
	chatUseCase        domain.ChatUseCase
	translationUseCase domain.TranslateUseCase
}

func NewAiController(l *slog.Logger, chatUseCase domain.ChatUseCase, translationUseCase domain.TranslateUseCase) *AiController {
	return &AiController{
		l:                  l,
		chatUseCase:        chatUseCase,
		translationUseCase: translationUseCase,
	}
}

func (c *AiController) Router(r chi.Router) {
	r.Post("/sentence", c.GenerateSentence)
	r.Post("/translate", c.Translate)
}

// GenerateSentence is an endpoint handler for generating a sentence containing submitted word.
func (c *AiController) GenerateSentence(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*30)
	defer cancel()

	var sentenceGenerationRequest domain.SentenceGenerationRequest
	if err := json.NewDecoder(r.Body).Decode(&sentenceGenerationRequest); err != nil {
		apiutil.Err(c.l, w, apiutil.ApiError{
			Status:  http.StatusBadRequest,
			Message: "Invalid sentence generation request payload",
			Cause:   err,
		})
		return
	}

	generatedSentence, err := c.chatUseCase.GenerateSentence(ctx, &sentenceGenerationRequest)
	if err != nil {
		if errors.Is(err, usecase.ErrValidation) {
			apiutil.Err(c.l, w, apiutil.ApiError{
				Status:  http.StatusBadRequest,
				Message: "Invalid request body",
				Cause:   err,
			})
		} else {
			apiutil.Err(c.l, w, err)
		}
		return
	}

	apiutil.Json(c.l, w, http.StatusOK, map[string]string{
		"sentence": generatedSentence,
	})
}

// Translate is an endpoint handler for translating words using DeepL.
func (c *AiController) Translate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*30)
	defer cancel()

	var body domain.TranslateRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apiutil.Err(c.l, w, apiutil.ApiError{
			Status:  http.StatusBadRequest,
			Message: "Invalid translation request payload",
			Cause:   err,
		})
		return
	}

	t, err := c.translationUseCase.Translate(ctx, &body)
	if err != nil {
		if errors.Is(err, usecase.ErrValidation) {
			apiutil.Err(c.l, w, apiutil.ApiError{
				Status:  http.StatusBadRequest,
				Message: "Invalid request body",
				Cause:   err,
			})
		} else {
			apiutil.Err(c.l, w, err)
		}
		return
	}

	apiutil.Json(c.l, w, http.StatusOK, map[string]string{
		"definition": t,
	})
}
