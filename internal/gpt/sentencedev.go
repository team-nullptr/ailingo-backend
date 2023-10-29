package gpt

import (
	"context"
	"fmt"
	"time"

	"ailingo/internal/domain"
)

type SentenceDevRepo struct{}

func NewSentenceDevRepo() domain.SentenceRepo {
	return &SentenceDevRepo{}
}

func (d SentenceDevRepo) GenerateSentence(ctx context.Context, word *domain.Definition) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(time.Second * 3):
		if word.Phrase == "fail" {
			return "", fmt.Errorf("%w: Example failure reason", ErrGenerationUnsuccessful)
		}
		return "What’s more amazing than a talking dog? A spelling bee.", nil
	}
}
