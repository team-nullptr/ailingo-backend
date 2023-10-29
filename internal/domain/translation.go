package domain

import "context"

type TranslationRepo interface {
	Translate(ctx context.Context, phrase string) (string, error)
}

type TranslateUseCase interface {
	Translate(ctx context.Context, phrase string) (string, error)
}
