package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	svix "github.com/svix/svix-webhooks/go"

	"ailingo/config"
	"ailingo/internal/domain"
	"ailingo/pkg/apiutil"
)

type ClerkWebhook struct {
	l           *slog.Logger
	cfg         *config.Config
	userUseCase domain.UserUseCase

	wh *svix.Webhook
}

func NewClerkWebhook(l *slog.Logger, cfg *config.Config, userUseCase domain.UserUseCase) (*ClerkWebhook, error) {
	wh, err := svix.NewWebhook(cfg.Services.ClerkWebhookSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to init the webhook: %w", err)
	}

	return &ClerkWebhook{
		l:           l,
		cfg:         cfg,
		userUseCase: userUseCase,
		wh:          wh,
	}, nil
}

type event struct {
	Data json.RawMessage `json:"data"`
	Type string          `json:"type"`
}

func (wh *ClerkWebhook) Webhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		apiutil.Err(wh.l, w, apiutil.ApiError{
			Status: http.StatusBadRequest,
			Cause:  err,
		})
		return
	}

	if err := wh.wh.Verify(body, r.Header); err != nil {
		apiutil.Err(wh.l, w, apiutil.ApiError{
			Status: http.StatusBadRequest,
			Cause:  err,
		})
		return
	}

	var ev event
	if err := json.Unmarshal(body, &ev); err != nil {
		apiutil.Err(wh.l, w, apiutil.ApiError{
			Status: http.StatusBadRequest,
			Cause:  err,
		})
		return
	}

	switch ev.Type {
	case "user.created":
		wh.userCreatedEventHandler(ctx, w, ev.Data)
	case "user.updated":
		wh.userUpdatedEventHandler(ctx, w, ev.Data)
	case "user.deleted":
		wh.userDeletedEventHandler(ctx, w, ev.Data)
	}
}

type userCreatedEventData struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	ImageURL string `json:"image_url"`
}

func (wh *ClerkWebhook) userCreatedEventHandler(ctx context.Context, w http.ResponseWriter, eventDataRaw json.RawMessage) {
	var eventData userCreatedEventData
	if err := json.Unmarshal(eventDataRaw, &eventData); err != nil {
		apiutil.Err(wh.l, w, apiutil.ApiError{
			Status: http.StatusBadRequest,
			Cause:  err,
		})
		return
	}

	if err := wh.userUseCase.Insert(ctx, &domain.InsertUserData{
		Id:       eventData.Id,
		Username: eventData.Username,
		ImageURL: eventData.ImageURL,
	}); err != nil {
		apiutil.Err(wh.l, w, err)
		return
	}

	apiutil.Empty(w, http.StatusCreated)
}

type userUpdatedEventData userCreatedEventData

func (wh *ClerkWebhook) userUpdatedEventHandler(ctx context.Context, w http.ResponseWriter, eventDataRaw json.RawMessage) {
	var eventData userUpdatedEventData
	if err := json.Unmarshal(eventDataRaw, &eventData); err != nil {
		apiutil.Err(wh.l, w, apiutil.ApiError{
			Status: http.StatusBadRequest,
			Cause:  err,
		})
		return
	}

	if err := wh.userUseCase.Update(ctx, &domain.UpdateUserData{
		Id:       eventData.Id,
		Username: eventData.Username,
		ImageURL: eventData.ImageURL,
	}); err != nil {
		apiutil.Err(wh.l, w, err)
		return
	}

	apiutil.Empty(w, http.StatusOK)
}

type userDeletedEventData struct {
	Id string `json:"id"`
}

func (wh *ClerkWebhook) userDeletedEventHandler(ctx context.Context, w http.ResponseWriter, eventDataRaw json.RawMessage) {
	var eventData userDeletedEventData
	if err := json.Unmarshal(eventDataRaw, &eventData); err != nil {
		apiutil.Err(wh.l, w, apiutil.ApiError{
			Status: http.StatusBadRequest,
			Cause:  err,
		})
		return
	}

	if err := wh.userUseCase.Delete(ctx, eventData.Id); err != nil {
		apiutil.Err(wh.l, w, err)
		return
	}

	apiutil.Empty(w, http.StatusOK)
}
