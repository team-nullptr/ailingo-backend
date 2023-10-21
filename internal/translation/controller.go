package translation

import (
	"context"
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

// Translate is an endpoint handler for translating words using DeepL.
func (c *Controller) Translate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*30)
	defer cancel()

	// TODO: Read body instead of query params
	phrase := r.URL.Query().Get("phrase")
	if phrase == "" {
		apiutil.Err(w, http.StatusBadRequest, errors.New("phrase query missing"))
		return
	}

	t, err := c.translator.Translate(ctx, phrase)
	if err != nil {
		apiutil.Err(w, http.StatusInternalServerError, err)
		return
	}

	apiutil.Json(w, http.StatusOK, map[string]string{
		"translation": t,
	})
}
