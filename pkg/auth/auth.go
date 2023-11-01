package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/clerkinc/clerk-sdk-go/clerk"

	apiutil2 "ailingo/pkg/apiutil"
)

var (
	ErrNoClaims = errors.New("no claims found in context")
)

type UserService struct {
	client clerk.Client
	logger *slog.Logger
}

func NewUserService(logger *slog.Logger, client clerk.Client) *UserService {
	return &UserService{
		client: client,
		logger: logger,
	}
}

// GetUserFromContext lookups the user whose claims were found in the context.
// In order for this function to work, the WithClaims middleware must be applied.
func (us *UserService) GetUserFromContext(ctx context.Context) (*clerk.User, error) {
	claims, ok := clerk.SessionFromContext(ctx)
	if !ok {
		return nil, ErrNoClaims
	}

	user, err := us.client.Users().Read(claims.Subject)
	if err != nil {
		return nil, fmt.Errorf("failed to read the user: %w", err)
	}

	return user, nil
}

func (us *UserService) GetUserById(userID string) (*clerk.User, error) {
	user, err := us.client.Users().Read(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to read the user: %w", err)
	}

	return user, nil
}

// WithClaims retrieves user auth token from the request and appends claims to the context.
func WithClaims(logger *slog.Logger, client clerk.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authToken := getAuthToken(r)
			claims, err := client.VerifyToken(authToken)
			if err != nil {
				apiutil2.Err(logger, w, apiutil2.ApiError{
					Status: http.StatusUnauthorized,
					Cause:  err,
				})
				return
			}

			ctx := context.WithValue(r.Context(), clerk.ActiveSessionClaims, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getAuthToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	return strings.TrimPrefix(strings.TrimSpace(authHeader), "Bearer ")
}
