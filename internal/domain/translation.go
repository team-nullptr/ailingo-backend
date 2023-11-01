package domain

import "context"

// TranslateRequest represents a translate request payload.
type TranslateRequest struct {
	Phrase string `json:"phrase" validate:"required,max=256"`
}

// TranslateRepo describes methods required by TranslateRepo implementations.
type TranslateRepo interface {
	Translate(ctx context.Context, phrase string) (string, error)
}

// TranslateUseCase describes methods required by TranslateUseCase implementation.
type TranslateUseCase interface {
	// Translate translates the given phrase into polish.
	Translate(ctx context.Context, translateRequest *TranslateRequest) (string, error)
}
