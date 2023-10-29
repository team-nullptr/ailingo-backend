package usecase

import (
	"context"
	"errors"
	"time"

	"ailingo/internal/domain"
)

type TranslateDevUseCase struct{}

func NewTranslateDevUseCase() domain.TranslateUseCase {
	return &TranslateDevUseCase{}
}

func (d *TranslateDevUseCase) Translate(ctx context.Context, phrase string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(time.Second * 3):
		if phrase == "fail" {
			return "", errors.New("something bad happened")
		}
		return "development", nil
	}
}
