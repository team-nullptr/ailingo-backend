package sentence

import (
	"ailingo/internal/models"
	"context"
	"fmt"
	"time"
)

type DevRepo struct{}

func NewDevRepo() Repo {
	return &DevRepo{}
}

var (
	errExampleFailure = "this is an example reason of failure"
)

func (d DevRepo) GenerateSentence(ctx context.Context, word models.Word) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(time.Second * 3):
		if word.Phrase == "fail" {
			return "", fmt.Errorf("%w: %s", ErrGenerationUnsuccessful, errExampleFailure)
		}
		return "What’s more amazing than a talking dog? A spelling bee.", nil
	}
}
