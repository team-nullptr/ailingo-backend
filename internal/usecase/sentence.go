package usecase

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"

	"ailingo/internal/domain"
)

// ChatUseCase expose features related with OpenAI's chat completion API.
type ChatUseCase struct {
	chatRepo domain.SentenceRepo
	validate *validator.Validate
}

func NewChatUseCase(repo domain.SentenceRepo, validate *validator.Validate) domain.ChatUseCase {
	return &ChatUseCase{
		chatRepo: repo,
		validate: validate,
	}
}

// GenerateSentence requests a new chat completion with Sentence Generator Persona prompt.
func (uc *ChatUseCase) GenerateSentence(ctx context.Context, sentenceGenerationRequest *domain.SentenceGenerationRequest) (string, error) {
	if err := uc.validate.Struct(sentenceGenerationRequest); err != nil {
		return "", fmt.Errorf("%w: %w", ErrValidation, err)
	}

	return uc.chatRepo.GenerateSentence(ctx, sentenceGenerationRequest)
}
