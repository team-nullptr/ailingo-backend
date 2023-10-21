package chat

import (
	"context"
	"time"

	"ailingo/internal/models"
)

// DevService can be used when the backend is running in development mode.
// It does not make real requests to OpenAI API. Returns a fake response instead
// to save money.
type DevService struct{}

func NewDevService() Service {
	return &DevService{}
}

// GenerateSentence generates a fake example sentence.
// If the phrase is equal to fail an example unsuccessful result will be returned.
func (ds *DevService) GenerateSentence(ctx context.Context, word models.Word) (*GenerationResult, error) {
	time.Sleep(time.Second * 3)

	if word.Phrase == "fail" {
		return &GenerationResult{
			Success: false,
			Reason:  "This is an example reason why I have failed",
		}, nil
	}

	return &GenerationResult{
		Success:  true,
		Sentence: "Why have you been farting around here all this time?",
	}, nil
}
