package studyset

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"ailingo/pkg/apiutil"
)

// Controller exposes handlers for study sets management API.
type Controller struct {
	studySetSvc *Service
	logger      *slog.Logger
}

func NewController(logger *slog.Logger, studySetSvc *Service) *Controller {
	return &Controller{
		logger:      logger,
		studySetSvc: studySetSvc,
	}
}

// Attach attaches controller to the given mux.
func (c *Controller) Attach(m *chi.Mux, path string) {
	m.Route(path, func(r chi.Router) {
		r.Post("/", c.Create)
	})
}

// Create is an endpoint handler for creating study sets.
func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	var body StudySetCreate
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apiutil.Err(c.logger, w, http.StatusBadRequest, err)
		return
	}

	created, err := c.studySetSvc.Create(&body)
	if err != nil {
		// TODO: Return a proper validation error
		if errors.Is(err, ErrValidation) {
			apiutil.Err(c.logger, w, http.StatusBadRequest, err)
		} else {
			apiutil.Err(c.logger, w, http.StatusInternalServerError, err)
		}
		return
	}

	apiutil.Json(c.logger, w, http.StatusCreated, created)
}
