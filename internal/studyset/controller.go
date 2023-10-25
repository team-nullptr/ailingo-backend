package studyset

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/go-chi/chi/v5"

	"ailingo/pkg/apiutil"
)

// Controller exposes handlers for study sets management API.
type Controller struct {
	logger          *slog.Logger
	studySetUseCase UseCase
}

func New(logger *slog.Logger, studySetUseCase UseCase) *Controller {
	return &Controller{
		logger:          logger,
		studySetUseCase: studySetUseCase,
	}
}

// Create is an endpoint handler for creating study sets.
func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	claims, ok := clerk.SessionFromContext(ctx)
	if !ok {
		apiutil.Err(c.logger, w, http.StatusUnauthorized, nil)
	}

	var body InsertStudySetData
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apiutil.Err(c.logger, w, http.StatusBadRequest, err)
		return
	}

	body.AuthorId = claims.Subject

	createdStudySet, err := c.studySetUseCase.Create(&body)
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

func (c *Controller) GetById(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		apiutil.Err(c.logger, w, http.StatusBadRequest, errors.New("invalid study set id"))
		return
	}

	studySet, err := c.studySetUseCase.GetById(id)
	if err != nil {
		apiutil.Err(c.logger, w, http.StatusInternalServerError, err)
		return
	}

	apiutil.Json(c.logger, w, http.StatusOK, studySet)
}
