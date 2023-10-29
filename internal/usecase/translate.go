package usecase

import (
	"context"

	"ailingo/internal/domain"
)

type TranslateUseCase struct {
	translationRepo domain.TranslationRepo
}

func NewTranslateUseCase(translationRepo domain.TranslationRepo) domain.TranslateUseCase {
	return &TranslateUseCase{
		translationRepo: translationRepo,
	}
}

func (uc *TranslateUseCase) Translate(ctx context.Context, phrase string) (string, error) {
	return uc.translationRepo.Translate(ctx, phrase)
}
