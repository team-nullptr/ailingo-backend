package translation

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"ailingo/pkg/apiutil"
)

// Controller exposes handlers for translation API.
type Controller struct {
	translator Translator
}

func NewController(translator Translator) *Controller {
	return &Controller{
		translator: translator,
	}
}

// Attach attaches controller to the given mux.
func (c *Controller) Attach(m *chi.Mux, path string) {
	m.Route(path, func(r chi.Router) {
		r.Post("/", c.Translate)
	})
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
		apiutil.Err(w, http.StatusBadRequest, errors.New("unprocessable request body"))
		return
	}

	t, err := c.translator.Translate(ctx, body.Phrase)
	if err != nil {
		apiutil.Err(w, http.StatusInternalServerError, err)
		return
	}

	apiutil.Json(w, http.StatusOK, map[string]string{
		"definition": t,
	})
}
