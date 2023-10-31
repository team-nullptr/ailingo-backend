package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"

	"ailingo/internal/domain"
)

type TranslateUseCase struct {
	translateRepo domain.TranslateRepo
	validate      *validator.Validate
}

func NewTranslateUseCase(translateRepo domain.TranslateRepo, validate *validator.Validate) domain.TranslateUseCase {
	return &TranslateUseCase{
		translateRepo: translateRepo,
		validate:      validate,
	}
}

func (uc *TranslateUseCase) Translate(ctx context.Context, translateRequest *domain.TranslateRequest) (string, error) {
	if err := uc.validate.Struct(translateRequest); err != nil {
		return "", fmt.Errorf("%w: %w", ErrValidation, err)
	}

	return uc.translateRepo.Translate(ctx, translateRequest.Phrase)
}

type TranslateDevUseCase struct{}

func NewTranslateDevUseCase() domain.TranslateUseCase {
	return &TranslateDevUseCase{}
}

func (d *TranslateDevUseCase) Translate(ctx context.Context, translateRequest *domain.TranslateRequest) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(time.Second * 3):
		if translateRequest.Phrase == "fail" {
			return "", errors.New("something bad happened")
		} else if translateRequest.Phrase == "invalid" {
			return "", ErrValidation
		}
		return "development", nil
	}
}
