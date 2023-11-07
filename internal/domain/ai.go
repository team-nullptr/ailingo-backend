package domain

import "context"

// SentenceGenerationRequest represents sentence generation request payload.
type SentenceGenerationRequest struct {
	Phrase  string `json:"phrase" validate:"required,max=256"`
	Meaning string `json:"meaning" validate:"required,max=256"`
}

type SetGenerationRequest struct {
	Name string `json:"name" validate:"required,max=256"`
}

// AiService describes methods required by SentenceRepo implementation.
type AiService interface {
	GenerateSentence(ctx context.Context, req *SentenceGenerationRequest) (string, error)
	GenerateDefinitions(ctx context.Context, req *SetGenerationRequest) ([]*InsertDefinitionData, error)
}

// ChatUseCase describes methods required by ChatUseCase implementation.
type ChatUseCase interface {
	GenerateSentence(ctx context.Context, req *SentenceGenerationRequest) (string, error)
}
