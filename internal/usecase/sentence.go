package usecase

import (
	"context"

	"ailingo/internal/domain"
)

// ChatUseCase expose features related with OpenAI's chat completion API.
type ChatUseCase struct {
	chatRepo domain.SentenceRepo
}

func NewChatUseCase(repo domain.SentenceRepo) domain.ChatUseCase {
	return &ChatUseCase{
		chatRepo: repo,
	}
}

// GenerateSentence requests a new chat completion with Sentence Generator Persona prompt.
func (uc *ChatUseCase) GenerateSentence(ctx context.Context, word *domain.Definition) (string, error) {
	return uc.chatRepo.GenerateSentence(ctx, word)
}
