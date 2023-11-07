package usecase

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"

	"ailingo/internal/domain"
)

// ChatUseCase expose features related with OpenAI's chat completion API.
type ChatUseCase struct {
	aiService domain.AiService
	validate  *validator.Validate
}

func NewChatUseCase(aiRepo domain.AiService, validate *validator.Validate) domain.ChatUseCase {
	return &ChatUseCase{
		aiService: aiRepo,
		validate:  validate,
	}
}

// GenerateSentence requests a new chat completion with Sentence Generator Persona.
func (uc *ChatUseCase) GenerateSentence(ctx context.Context, req *domain.SentenceGenerationRequest) (string, error) {
	if err := uc.validate.Struct(req); err != nil {
		return "", fmt.Errorf("%w: %w", ErrValidation, err)
	}
	return uc.aiService.GenerateSentence(ctx, req)
}
