package translation

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"ailingo/pkg/apiutil"
)

// Controller exposes handlers for translation API.
type Controller struct {
	translator Translator
	logger     *slog.Logger
}

func NewController(logger *slog.Logger, translator Translator) *Controller {
	return &Controller{
		translator: translator,
		logger:     logger,
	}
}

type TranslatePayload struct {
	Phrase string `json:"phrase"`
}

// Translate is an endpoint handler for translating words using DeepL.
func (c *Controller) Translate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*30)
	defer cancel()

	var body TranslatePayload
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apiutil.Err(c.logger, w, http.StatusBadRequest, errors.New("unprocessable request body"))
		return
	}

	t, err := c.translator.Translate(ctx, body.Phrase)
	if err != nil {
		apiutil.Err(c.logger, w, http.StatusInternalServerError, err)
		return
	}

	apiutil.Json(c.logger, w, http.StatusOK, map[string]string{
		"definition": t,
	})
}
