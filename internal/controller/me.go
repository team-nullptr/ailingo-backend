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
	l              *slog.Logger
	profileUseCase domain.ProfileUseCase
	userService    *auth.UserService
}

func NewMeController(l *slog.Logger, accountUseCase domain.ProfileUseCase, userService *auth.UserService) *MeController {
	return &MeController{
		l:              l,
		profileUseCase: accountUseCase,
		userService:    userService,
	}
}

func (c *MeController) Router(r chi.Router) {
	r.Get("/study-sets/created", c.GetCreated)
	r.Get("/study-sets/starred", c.GetStarred)
	r.Post("/study-sets/starred", c.Star)
	r.Delete("/study-sets/starred/{studySetID}", c.Instar)
}

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
