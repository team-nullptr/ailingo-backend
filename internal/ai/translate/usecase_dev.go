package translate

import (
	"context"
	"time"
)

type TranslationDevUseCase struct{}

func NewTranslationDev() TranslationUseCase {
	return &TranslationDevUseCase{}
}

func (dt *TranslationDevUseCase) Translate(ctx context.Context, text string) (string, error) {
	time.Sleep(time.Second * 2)
	return "development", nil
}
