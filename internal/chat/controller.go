package chat

import (
	"ailingo/internal/models"
	"ailingo/pkg/apiutil"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Controller struct {
	ss *SentenceService
}

func NewController(ss *SentenceService) *Controller {
	return &Controller{
		ss: ss,
	}
}

func (c *Controller) Attach(m *chi.Mux, path string) {
	m.Route(path, func(r chi.Router) {
		r.Post("/sentence", c.GenerateSentence)
	})
}

// GenerateSentence is an endpoint handler that generates example sentences for given word.
func (c *Controller) GenerateSentence(w http.ResponseWriter, r *http.Request) {
	var word models.Word

	if err := json.NewDecoder(r.Body).Decode(&word); err != nil {
		apiutil.Err(w, http.StatusBadRequest, err)
		return
	}

	result, err := c.ss.GenerateSentence(word)
	if err != nil {
		apiutil.Err(w, http.StatusInternalServerError, err)
		return
	}

	apiutil.Json(w, http.StatusOK, result)
}
