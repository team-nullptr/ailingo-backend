package controller

import (
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

type TaskController struct {
	l           *slog.Logger
	userService *auth.UserService
	taskUseCase domain.TaskUseCase
}

func NewTaskController(l *slog.Logger, userService *auth.UserService, taskUseCase domain.TaskUseCase) *TaskController {
	return &TaskController{
		l:           l,
		userService: userService,
		taskUseCase: taskUseCase,
	}
}

func (c *TaskController) Router(r chi.Router) {
	r.Get("/{taskID}", c.Get)
}

func (c *TaskController) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	taskID, err := strconv.ParseInt(chi.URLParam(r, "taskID"), 10, 64)
	if err != nil {
		apiutil.Err(c.l, w, &apiutil.ApiError{
			Status:  http.StatusBadRequest,
			Message: "Invalid task ID",
		})
		return
	}

	task, err := c.taskUseCase.Get(ctx, taskID)
	if err != nil {
		var errNotFound *usecase.ErrNotFound
		if errors.As(err, &errNotFound) {
			apiutil.Err(c.l, w, &apiutil.ApiError{
				Status:  http.StatusNotFound,
				Message: errNotFound.Error(),
			})
		} else {
			apiutil.Err(c.l, w, err)
		}
		return
	}

	apiutil.Json(c.l, w, http.StatusOK, task)
}
