package sentence

import (
	"context"
	_ "embed"

	"ailingo/internal/models"
)

type ChatUseCase interface {
	GenerateSentence(ctx context.Context, word models.Definition) (string, error)
}

// DefaultChatUseCase expose features related with OpenAI's chat completion API.
type DefaultChatUseCase struct {
	chatRepo Repo
}

func NewChatUseCase(repo Repo) ChatUseCase {
	return &DefaultChatUseCase{
		chatRepo: repo,
	}
}

// GenerateSentence requests a new chat completion with Sentence Generator Persona prompt.
func (uc *DefaultChatUseCase) GenerateSentence(ctx context.Context, word models.Definition) (string, error) {
	return uc.chatRepo.GenerateSentence(ctx, word)
}
