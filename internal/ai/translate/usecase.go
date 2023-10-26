package translate

import "context"

type TranslationUseCase interface {
	Translate(ctx context.Context, phrase string) (string, error)
}

type DefaultTranslationUseCase struct {
	translationRepo Repo
}

func NewTranslationUseCase(translationRepo Repo) TranslationUseCase {
	return &DefaultTranslationUseCase{
		translationRepo: translationRepo,
	}
}

func (uc DefaultTranslationUseCase) Translate(ctx context.Context, phrase string) (string, error) {
	return uc.translationRepo.TranslatePhrase(ctx, phrase)
}
