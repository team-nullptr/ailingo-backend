package controller

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"ailingo/internal/domain"
	"ailingo/internal/usecase"
	"ailingo/pkg/apiutil"
	"ailingo/pkg/auth"
)

type MeController struct {
	l                   *slog.Logger
	profileUseCase      domain.ProfileUseCase
	studySessionUseCase domain.StudySessionUseCase
	userService         *auth.UserService
}

func NewMeController(l *slog.Logger, accountUseCase domain.ProfileUseCase, studySessionUseCase domain.StudySessionUseCase, userService *auth.UserService) *MeController {
	return &MeController{
		l:                   l,
		profileUseCase:      accountUseCase,
		studySessionUseCase: studySessionUseCase,
		userService:         userService,
	}
}

func (c *MeController) Router(r chi.Router) {
	r.Get("/study-sets/created", c.GetCreated)
	r.Get("/study-sets/starred", c.GetStarred)
	r.Post("/study-sets/starred", c.Star)
	r.Delete("/study-sets/starred/{studySetID}", c.Instar)

	r.Get("/study-sessions", c.GetRecentStudySessions)
	r.Get("/study-sessions/{studySetID}", c.GetStudySessionForStudySet)
	r.Patch("/study-sessions/{studySetID}", c.RefreshStudySession)
}

// GetCreated is an endpoint handler for getting all created study sets.
func (c *MeController) GetCreated(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	studySets, err := c.profileUseCase.GetCreatedStudySets(ctx, user.ID)
	if err != nil {
		apiutil.Err(c.l, w, err)
		return
	}

	apiutil.Json(c.l, w, http.StatusOK, studySets)
}

// GetStarred is an endpoint handler for getting starred study sets.
func (c *MeController) GetStarred(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	studySets, err := c.profileUseCase.GetStarredStudySets(ctx, user.ID)
	if err != nil {
		apiutil.Err(c.l, w, err)
		return
	}

	apiutil.Json(c.l, w, http.StatusOK, studySets)
}

type starPayload struct {
	Id int64 `json:"id"`
}

// Star is an endpoint handler for adding the given study set to the starred study sets list.
func (c *MeController) Star(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	var body starPayload
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apiutil.Err(c.l, w, apiutil.ApiError{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Cause:   err,
		})
		return
	}

	if err := c.profileUseCase.StarStudySet(ctx, user.ID, body.Id); err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			apiutil.Empty(w, http.StatusNotFound)
		} else if errors.Is(err, usecase.ErrAlreadyStarred) {
			apiutil.Empty(w, http.StatusBadRequest)
		} else {
			apiutil.Err(c.l, w, err)
		}
		return
	}

	apiutil.Empty(w, http.StatusCreated)
}

// Instar is an endpoint handler for removing the given study set from starred study sets list.
func (c *MeController) Instar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	if err := c.profileUseCase.InstarStudySet(ctx, user.ID, studySetID); err != nil {
		apiutil.Empty(w, http.StatusInternalServerError)
		return
	}

	apiutil.Empty(w, http.StatusOK)
}

// GetRecentStudySessions is an endpoint for getting recent study sessions for the authenticated user.
func (c *MeController) GetRecentStudySessions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	studySessions, err := c.studySessionUseCase.GetRecent(ctx, user.ID)
	if err != nil {
		apiutil.Err(c.l, w, err)
		return
	}

	apiutil.Json(c.l, w, http.StatusOK, studySessions)
}

// GetStudySessionForStudySet is an endpoint handler for getting study session for the given study set.
func (c *MeController) GetStudySessionForStudySet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	studySession, err := c.studySessionUseCase.GetForStudySet(ctx, user.ID, studySetID)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			apiutil.Empty(w, http.StatusNotFound)
			return
		} else {
			apiutil.Err(c.l, w, err)
		}
		return
	}

	apiutil.Json(c.l, w, http.StatusOK, studySession)
}

// RefreshStudySession is an endpoint handler used for updating user study session's last session timestamp
// or creating a new study session (if it does not exist yet).
func (c *MeController) RefreshStudySession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	if err := c.studySessionUseCase.Refresh(ctx, user.ID, studySetID); err != nil {
		if errors.Is(err, usecase.ErrNotFound) {
			apiutil.Empty(w, http.StatusNotFound)
		} else {
			apiutil.Err(c.l, w, err)
		}
		return
	}

	apiutil.Empty(w, http.StatusOK)
}
