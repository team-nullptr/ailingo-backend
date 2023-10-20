package translation

import (
	"ailingo/pkg/apiutil"
	"ailingo/pkg/deepl"
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

type Controller struct {
	deeplClient *deepl.Client
}

func NewController(deeplClient *deepl.Client) *Controller {
	return &Controller{
		deeplClient: deeplClient,
	}
}

func (c *Controller) Attach(m *chi.Mux, path string) {
	m.Route(path, func(r chi.Router) {
		r.Post("/", c.Translate)
	})
}

func (c *Controller) Translate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*30)
	defer cancel()

	// TODO: Read body instead of query params
	phrase := r.URL.Query().Get("phrase")
	if phrase == "" {
		apiutil.Err(w, http.StatusBadRequest, errors.New("phrase query missing"))
		return
	}

	t, err := c.deeplClient.Translate(ctx, phrase)
	if err != nil {
		apiutil.Err(w, http.StatusInternalServerError, err)
		return
	}

	apiutil.Json(w, http.StatusOK, map[string]string{
		"translation": t,
	})
}
