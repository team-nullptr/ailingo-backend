package auth

import (
	"ailingo/pkg/apiutil"
	"context"
	"errors"
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"net/http"
	"strings"
)

var (
	ErrNoClaims = errors.New("no claims found in context")
)

type UserService struct {
	client clerk.Client
}

func NewUserService(client clerk.Client) *UserService {
	return &UserService{
		client: client,
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
func WithClaims(client clerk.Client) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authToken := getAuthToken(r)
			claims, err := client.VerifyToken(authToken)
			if err != nil {
				apiutil.Err(w, http.StatusUnauthorized, nil)
				return
			}
			ctx := context.WithValue(r.Context(), clerk.ActiveSessionClaims, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

func getAuthToken(r *http.Request) string {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	return strings.TrimPrefix(authHeader, "Bearer ")
}
