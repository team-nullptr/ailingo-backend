package usecase

import "context"

type TranslationUseCase interface {
	Translate(ctx context.Context, text string) (string, error)
}
