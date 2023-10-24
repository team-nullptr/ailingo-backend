package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/clerkinc/clerk-sdk-go/clerk"

	"ailingo/pkg/apiutil"
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

// GetUser lookups the user whose claims were found in the context.
// In order for this function to work, the WithClaims middleware must be applied.
func (us UserService) GetUser(ctx context.Context) (*clerk.User, error) {
	claims, ok := clerk.SessionFromContext(ctx)
	if !ok {
		return nil, ErrNoClaims
	}
	user, err := us.client.Users().Read(claims.Subject)
	if err != nil {
		return nil, err
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
				apiutil.Err(logger, w, http.StatusUnauthorized, nil)
				return
			}

			ctx := context.WithValue(r.Context(), clerk.ActiveSessionClaims, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getAuthToken(r *http.Request) string {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	return strings.TrimPrefix(authHeader, "Bearer ")
}
