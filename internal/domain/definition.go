package domain

import "context"

type Definition struct {
	Id        int64    `json:"id"`
	Phrase    string   `json:"phrase"`
	Meaning   string   `json:"meaning"`
	Sentences []string `json:"sentences"`
}

// DefinitionRow is a "fake" type that exists in case we would want to store additional information
// in definition table.
type DefinitionRow Definition

func (r *DefinitionRow) Populate() *Definition {
	return &Definition{
		Id:        r.Id,
		Phrase:    r.Phrase,
		Meaning:   r.Meaning,
		Sentences: r.Sentences,
	}
}

type InsertDefinitionData struct {
	Phrase    string   `json:"phrase"`
	Meaning   string   `json:"meaning"`
	Sentences []string `json:"sentences"`
}

// DefinitionRepo describes methods required by DefinitionRepo implementation.
type DefinitionRepo interface {
	GetAllFor(ctx context.Context, studySetID int64) ([]*DefinitionRow, error)
	Insert(ctx context.Context, studySetID int64, insertData *InsertDefinitionData) error
}

// DefinitionUseCase describes methods required by DefinitionUseCase implementation.
type DefinitionUseCase interface {
	GetAllFor(ctx context.Context, studySetID int64) ([]*Definition, error)
	Create(ctx context.Context, userID string, studySetID int64, insertData *InsertDefinitionData) error
}
