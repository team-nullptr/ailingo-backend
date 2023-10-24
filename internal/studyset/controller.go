package studyset

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"

	"ailingo/pkg/apiutil"
)

// Controller exposes handlers for study sets management API.
type Controller struct {
	studySetService *Service
	logger          *slog.Logger
}

func NewController(logger *slog.Logger, studySetService *Service) *Controller {
	return &Controller{
		logger:          logger,
		studySetService: studySetService,
	}
}

// Create is an endpoint handler for creating study sets.
func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, ok := clerk.SessionFromContext(ctx)
	if !ok {
		apiutil.Err(c.logger, w, http.StatusUnauthorized, nil)
	}

	var body studySetCreateData
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apiutil.Err(c.logger, w, http.StatusBadRequest, err)
		return
	}

	body.AuthorId = claims.Subject

	createdStudySet, err := c.studySetService.Create(&body)
	if err != nil {
		// TODO: Return a proper validation error
		if errors.Is(err, ErrValidation) {
			apiutil.Err(c.logger, w, http.StatusBadRequest, err)
		} else {
			apiutil.Err(c.logger, w, http.StatusInternalServerError, err)
		}
		return
	}

	apiutil.Json(c.logger, w, http.StatusCreated, createdStudySet)
}
