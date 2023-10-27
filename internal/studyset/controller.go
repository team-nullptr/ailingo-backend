package studyset

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"ailingo/internal/apiutil"
	"ailingo/internal/auth"
)

// Controller exposes handlers for study sets management API.
type Controller struct {
	logger          *slog.Logger
	userService     *auth.UserService
	studySetUseCase UseCase
}

func NewController(logger *slog.Logger, userService *auth.UserService, studySetUseCase UseCase) *Controller {
	return &Controller{
		logger:          logger,
		userService:     userService,
		studySetUseCase: studySetUseCase,
	}
}

// GetAllSummary is an endpoint handler for getting a summary for all existing study sets.
func (c *Controller) GetAllSummary(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	studySets, err := c.studySetUseCase.GetAllSummary(ctx)
	if err != nil {
		apiutil.Err(c.logger, w, err)
		return
	}

	apiutil.Json(c.logger, w, http.StatusOK, studySets)
}

// GetById is an endpoint handler for getting full information about a specific study set.
func (c *Controller) GetById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	studySetID, err := strconv.ParseInt(chi.URLParam(r, "studySetID"), 10, 64)
	if err != nil {
		apiutil.Err(c.logger, w, apiutil.ApiError{
			Status:  http.StatusBadRequest,
			Message: "Invalid study set ID",
		})
		return
	}

	studySet, err := c.studySetUseCase.GetById(ctx, studySetID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			apiutil.Err(c.logger, w, apiutil.ApiError{
				Status:  http.StatusNotFound,
				Message: "This study set does not exist",
			})
		} else {
			apiutil.Err(c.logger, w, err)
		}
		return
	}

	apiutil.Json(c.logger, w, http.StatusOK, studySet)
}

// Create is an endpoint handler for creating a new study set.
func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	user, err := c.userService.GetUser(ctx)
	if err != nil {
		if errors.Is(err, auth.ErrNoClaims) {
			apiutil.Err(c.logger, w, apiutil.ApiError{
				Status: http.StatusUnauthorized,
				Cause:  err,
			})
		} else {
			apiutil.Err(c.logger, w, err)
		}
	}

	var insertData insertStudySetData
	if err := json.NewDecoder(r.Body).Decode(&insertData); err != nil {
		apiutil.Err(c.logger, w, apiutil.ApiError{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Cause:   err,
		})
		return
	}
	insertData.AuthorId = user.ID

	createdStudySet, err := c.studySetUseCase.Create(ctx, &insertData)
	if err != nil {
		if errors.Is(err, ErrValidation) {
			apiutil.Err(c.logger, w, apiutil.ApiError{
				Status:  http.StatusBadRequest,
				Message: "Invalid request body",
				Cause:   err,
			})
		} else {
			apiutil.Err(c.logger, w, err)
		}
		return
	}

	apiutil.Json(c.logger, w, http.StatusCreated, createdStudySet)
}

// Update is an endpoint for replacing data of existing study set.
func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	user, err := c.userService.GetUser(ctx)
	if err != nil {
		if errors.Is(err, auth.ErrNoClaims) {
			apiutil.Err(c.logger, w, apiutil.ApiError{
				Status:  http.StatusUnauthorized,
				Message: "Missing authorization",
				Cause:   err,
			})
		} else {
			apiutil.Err(c.logger, w, err)
		}
		return
	}

	studySetID, err := strconv.ParseInt(chi.URLParam(r, "studySetID"), 10, 64)
	if err != nil {
		apiutil.Err(c.logger, w, apiutil.ApiError{
			Status:  http.StatusBadRequest,
			Message: "Invalid study set ID",
		})
		return
	}

	var updateData updateStudySetData
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		apiutil.Err(c.logger, w, apiutil.ApiError{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Cause:   err,
		})
		return
	}

	if err := c.studySetUseCase.Update(ctx, studySetID, user.ID, &updateData); err != nil {
		if errors.Is(err, ErrNotFound) {
			apiutil.Err(c.logger, w, apiutil.ApiError{
				Status:  http.StatusNotFound,
				Message: "This study set does not exist",
			})
		} else if errors.Is(err, ErrForbidden) {
			apiutil.Err(c.logger, w, apiutil.ApiError{
				Status:  http.StatusForbidden,
				Message: "You don't have enough permission to update this study set",
				Cause:   err,
			})
		} else {
			apiutil.Err(c.logger, w, err)
		}
		return
	}

	apiutil.Empty(w, http.StatusOK)
}

// Delete is an endpoint for deleting a study set.
func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	user, err := c.userService.GetUser(ctx)
	if err != nil {
		if errors.Is(err, auth.ErrNoClaims) {
			apiutil.Err(c.logger, w, apiutil.ApiError{
				Status:  http.StatusUnauthorized,
				Message: "Missing authorization",
				Cause:   err,
			})
		} else {
			apiutil.Err(c.logger, w, err)
		}
		return
	}

	studySetID, err := strconv.ParseInt(chi.URLParam(r, "studySetID"), 10, 64)
	if err != nil {
		apiutil.Err(c.logger, w, apiutil.ApiError{
			Status:  http.StatusBadRequest,
			Message: "Invalid study set ID",
		})
		return
	}

	if err := c.studySetUseCase.Delete(ctx, studySetID, user.ID); err != nil {
		if errors.Is(err, ErrNotFound) {
			apiutil.Err(c.logger, w, apiutil.ApiError{
				Status:  http.StatusNotFound,
				Message: "This study set does not exist",
			})
		} else if errors.Is(err, ErrForbidden) {
			apiutil.Err(c.logger, w, apiutil.ApiError{
				Status:  http.StatusForbidden,
				Message: "You don't have enough permission to update this study set",
				Cause:   err,
			})
		} else {
			apiutil.Err(c.logger, w, err)
		}
		return
	}

	apiutil.Empty(w, http.StatusOK)
}
