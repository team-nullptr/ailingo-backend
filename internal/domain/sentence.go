package domain

import "context"

type SentenceRepo interface {
	GenerateSentence(ctx context.Context, word *Definition) (string, error)
}

type ChatUseCase interface {
	GenerateSentence(ctx context.Context, word *Definition) (string, error)
}
