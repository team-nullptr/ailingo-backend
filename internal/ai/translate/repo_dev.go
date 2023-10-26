package translate

import (
	"context"
	"errors"
	"time"
)

type DevRepo struct{}

func NewDevRepo() Repo {
	return &DevRepo{}
}

func (d DevRepo) TranslatePhrase(ctx context.Context, phrase string) (string, error) {
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
