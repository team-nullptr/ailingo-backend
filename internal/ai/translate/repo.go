package translate

import (
	"ailingo/pkg/deepl"
	"context"
)

type Repo interface {
	TranslatePhrase(ctx context.Context, phrase string) (string, error)
}

type DefaultRepo struct {
	translationClient *deepl.Client
}

func NewRepo(translationClient *deepl.Client) Repo {
	return &DefaultRepo{
		translationClient: translationClient,
	}
}

func (d DefaultRepo) TranslatePhrase(ctx context.Context, phrase string) (string, error) {
	return d.translationClient.Translate(ctx, phrase)
}
