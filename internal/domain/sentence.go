package domain

import "context"

// SentenceGenerationRequest represents sentence generation request payload.
type SentenceGenerationRequest struct {
	Phrase  string `json:"phrase" validate:"required,max=256"`
	Meaning string `json:"meaning" validate:"required,max=256"`
}

// SentenceRepo describes methods required by SentenceRepo implementation.
type SentenceRepo interface {
	GenerateSentence(ctx context.Context, sentenceGenerationRequest *SentenceGenerationRequest) (string, error)
}

// ChatUseCase describes methods required by ChatUseCase implementation.
type ChatUseCase interface {
	GenerateSentence(ctx context.Context, sentenceGenerationRequest *SentenceGenerationRequest) (string, error)
}
