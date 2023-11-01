package controller

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"ailingo/internal/domain"
	"ailingo/internal/usecase"
	"ailingo/pkg/apiutil"
	"ailingo/pkg/auth"
)

type StudySetController struct {
	l                 *slog.Logger
	userService       *auth.UserService
	studySetUseCase   domain.StudySetUseCase
	definitionUseCase domain.DefinitionUseCase
}

func NewStudySetController(l *slog.Logger, userService *auth.UserService, studySetUseCase domain.StudySetUseCase, definitionUseCase domain.DefinitionUseCase) *StudySetController {
	return &StudySetController{
		l:                 l,
		userService:       userService,
		studySetUseCase:   studySetUseCase,
		definitionUseCase: definitionUseCase,
	}
}

func (c *StudySetController) Router(withClaims func(next http.Handler) http.Handler) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/", c.GetAllSummary)
		r.With(withClaims).Post("/", c.Create)
		r.Get("/{studySetID}", c.GetById)
		r.With(withClaims).Put("/{studySetID}", c.Update)
		r.With(withClaims).Delete("/{studySetID}", c.Delete)

		// TODO: We could make a separate controller for /definitions endpoints
		r.Get("/{parentStudySetID}/definitions", c.GetDefinitions)
		r.With(withClaims).Post("/{parentStudySetID}/definitions", c.CreateDefinition)
		r.With(withClaims).Put("/{parentStudySetID}/definitions/{definitionID}", c.UpdateDefinition)
		r.With(withClaims).Delete("/{parentStudySetID}/definitions/{definitionID}", c.DeleteDefinition)
	}
}

// GetAllSummary is an endpoint handler for getting a summary for all existing study sets.
func (c *StudySetController) GetAllSummary(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	studySets, err := c.studySetUseCase.GetAll(ctx)
	if err != nil {
		apiutil.Err(c.l, w, err)
		return
	}

	apiutil.Json(c.l, w, http.StatusOK, studySets)
}

// GetById is an endpoint handler for getting full information about a specific study set.
func (c *StudySetController) GetById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	studySetID, err := strconv.ParseInt(chi.URLParam(r, "studySetID"), 10, 64)
	if err != nil {
		apiutil.Err(c.l, w, apiutil.ApiError{
			Status:  http.StatusBadRequest,
			Message: "Invalid study set ID",
		})
		return
	}

	studySet, err := c.studySetUseCase.GetById(ctx, studySetID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			apiutil.Err(c.l, w, apiutil.ApiError{
				Status: http.StatusNotFound,
			})
		} else {
			apiutil.Err(c.l, w, err)
		}
		return
	}

	apiutil.Json(c.l, w, http.StatusOK, studySet)
}

// Create is an endpoint handler for creating a new study set.
func (c *StudySetController) Create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	user, err := c.userService.GetUserFromContext(ctx)
	if err != nil {
		if errors.Is(err, auth.ErrNoClaims) {
			apiutil.Err(c.l, w, apiutil.ApiError{
				Status: http.StatusUnauthorized,
				Cause:  err,
			})
		} else {
			apiutil.Err(c.l, w, err)
		}
	}

	var insertData domain.InsertStudySetData
	if err := json.NewDecoder(r.Body).Decode(&insertData); err != nil {
		apiutil.Err(c.l, w, apiutil.ApiError{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Cause:   err,
		})
		return
	}

	// TODO: Logic leaking to the controller
	insertData.AuthorId = user.ID

	createdStudySet, err := c.studySetUseCase.Create(ctx, &insertData)
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

	apiutil.Json(c.l, w, http.StatusCreated, createdStudySet)
}

// Update is an endpoint for replacing data of existing study set.
func (c *StudySetController) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	user, err := c.userService.GetUserFromContext(ctx)
	if err != nil {
		if errors.Is(err, auth.ErrNoClaims) {
			apiutil.Err(c.l, w, apiutil.ApiError{
				Status:  http.StatusUnauthorized,
				Message: "Missing authorization",
				Cause:   err,
			})
		} else {
			apiutil.Err(c.l, w, err)
		}
		return
	}

	studySetID, err := strconv.ParseInt(chi.URLParam(r, "studySetID"), 10, 64)
	if err != nil {
		apiutil.Err(c.l, w, apiutil.ApiError{
			Status:  http.StatusBadRequest,
			Message: "Invalid study set ID",
		})
		return
	}

	var updateData domain.UpdateStudySetData
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		apiutil.Err(c.l, w, apiutil.ApiError{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Cause:   err,
		})
		return
	}

	if err := c.studySetUseCase.Update(ctx, studySetID, user.ID, &updateData); err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			apiutil.Err(c.l, w, apiutil.ApiError{
				Status: http.StatusNotFound,
			})
		} else if errors.Is(err, usecase.ErrForbidden) {
			apiutil.Err(c.l, w, apiutil.ApiError{
				Status:  http.StatusForbidden,
				Message: "You don't have enough permission to update this study set",
				Cause:   err,
			})
		} else {
			apiutil.Err(c.l, w, err)
		}
		return
	}

	apiutil.Empty(w, http.StatusOK)
}

// Delete is an endpoint for deleting a study set.
func (c *StudySetController) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	user, err := c.userService.GetUserFromContext(ctx)
	if err != nil {
		if errors.Is(err, auth.ErrNoClaims) {
			apiutil.Err(c.l, w, apiutil.ApiError{
				Status:  http.StatusUnauthorized,
				Message: "Missing authorization",
				Cause:   err,
			})
		} else {
			apiutil.Err(c.l, w, err)
		}
		return
	}

	studySetID, err := strconv.ParseInt(chi.URLParam(r, "studySetID"), 10, 64)
	if err != nil {
		apiutil.Err(c.l, w, apiutil.ApiError{
			Status:  http.StatusBadRequest,
			Message: "Invalid study set ID",
		})
		return
	}

	if err := c.studySetUseCase.Delete(ctx, studySetID, user.ID); err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			apiutil.Err(c.l, w, apiutil.ApiError{
				Status: http.StatusNotFound,
			})
		} else if errors.Is(err, usecase.ErrForbidden) {
			apiutil.Err(c.l, w, apiutil.ApiError{
				Status:  http.StatusForbidden,
				Message: "You don't have enough permission to update this study set",
				Cause:   err,
			})
		} else {
			apiutil.Err(c.l, w, err)
		}
		return
	}

	apiutil.Empty(w, http.StatusOK)
}

func (c *StudySetController) GetDefinitions(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	parentStudySetID, err := strconv.ParseInt(chi.URLParam(r, "parentStudySetID"), 10, 64)
	if err != nil {
		apiutil.Err(c.l, w, apiutil.ApiError{
			Status:  http.StatusBadRequest,
			Message: "Invalid study set ID",
		})
		return
	}

	definitions, err := c.definitionUseCase.GetAllFor(ctx, parentStudySetID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			apiutil.Err(c.l, w, apiutil.ApiError{
				Status: http.StatusNotFound,
			})
		} else {
			apiutil.Err(c.l, w, err)
		}
		return
	}

	apiutil.Json(c.l, w, http.StatusOK, definitions)
}

func (c *StudySetController) CreateDefinition(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	user, err := c.userService.GetUserFromContext(ctx)
	if err != nil {
		if errors.Is(err, auth.ErrNoClaims) {
			apiutil.Err(c.l, w, apiutil.ApiError{
				Status: http.StatusUnauthorized,
				Cause:  err,
			})
		} else {
			apiutil.Err(c.l, w, err)
		}
	}

	parentStudySetID, err := strconv.ParseInt(chi.URLParam(r, "parentStudySetID"), 10, 64)
	if err != nil {
		apiutil.Err(c.l, w, apiutil.ApiError{
			Status:  http.StatusBadRequest,
			Message: "Invalid study set ID",
		})
		return
	}

	var insertData domain.InsertDefinitionData
	if err := json.NewDecoder(r.Body).Decode(&insertData); err != nil {
		apiutil.Err(c.l, w, apiutil.ApiError{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Cause:   err,
		})
		return
	}

	if err := c.definitionUseCase.Create(ctx, user.ID, parentStudySetID, &insertData); err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			apiutil.Err(c.l, w, apiutil.ApiError{
				Status: http.StatusNotFound,
			})
		} else {
			apiutil.Err(c.l, w, err)
		}
		return
	}

	apiutil.Empty(w, http.StatusCreated)
}

func (c *StudySetController) UpdateDefinition(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	user, err := c.userService.GetUserFromContext(ctx)
	if err != nil {
		if errors.Is(err, auth.ErrNoClaims) {
			apiutil.Err(c.l, w, apiutil.ApiError{
				Status: http.StatusUnauthorized,
				Cause:  err,
			})
		} else {
			apiutil.Err(c.l, w, err)
		}
	}

	parentStudySetID, err := strconv.ParseInt(chi.URLParam(r, "parentStudySetID"), 10, 64)
	if err != nil {
		apiutil.Err(c.l, w, apiutil.ApiError{
			Status:  http.StatusBadRequest,
			Message: "Invalid study set ID",
		})
		return
	}

	definitionID, err := strconv.ParseInt(chi.URLParam(r, "definitionID"), 10, 64)
	if err != nil {
		apiutil.Err(c.l, w, apiutil.ApiError{
			Status:  http.StatusBadRequest,
			Message: "Invalid definition ID",
		})
		return
	}

	var updateData domain.UpdateDefinitionData
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		apiutil.Err(c.l, w, apiutil.ApiError{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Cause:   err,
		})
		return
	}

	if err := c.definitionUseCase.Update(ctx, user.ID, parentStudySetID, definitionID, &updateData); err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			apiutil.Err(c.l, w, apiutil.ApiError{
				Status: http.StatusNotFound,
			})
		} else {
			apiutil.Err(c.l, w, err)
		}
		return
	}

	apiutil.Empty(w, http.StatusOK)
}

func (c *StudySetController) DeleteDefinition(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	user, err := c.userService.GetUserFromContext(ctx)
	if err != nil {
		if errors.Is(err, auth.ErrNoClaims) {
			apiutil.Err(c.l, w, apiutil.ApiError{
				Status: http.StatusUnauthorized,
				Cause:  err,
			})
		} else {
			apiutil.Err(c.l, w, err)
		}
	}

	parentStudySetID, err := strconv.ParseInt(chi.URLParam(r, "parentStudySetID"), 10, 64)
	if err != nil {
		apiutil.Err(c.l, w, apiutil.ApiError{
			Status:  http.StatusBadRequest,
			Message: "Invalid study set ID",
		})
		return
	}

	definitionID, err := strconv.ParseInt(chi.URLParam(r, "definitionID"), 10, 64)
	if err != nil {
		apiutil.Err(c.l, w, apiutil.ApiError{
			Status:  http.StatusBadRequest,
			Message: "Invalid definition set ID",
		})
		return
	}

	if err := c.definitionUseCase.Delete(ctx, user.ID, parentStudySetID, definitionID); err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			apiutil.Err(c.l, w, apiutil.ApiError{
				Status: http.StatusNotFound,
			})
		} else {
			apiutil.Err(c.l, w, err)
		}
		return
	}

	apiutil.Empty(w, http.StatusOK)
}
