package domain

import "context"

type Definition struct {
	Id        int64    `json:"id"`
	Phrase    string   `json:"phrase"`
	Meaning   string   `json:"meaning"`
	Sentences []string `json:"sentences"`
}

type InsertDefinitionData struct {
	Phrase    string   `json:"phrase"`
	Meaning   string   `json:"meaning"`
	Sentences []string `json:"sentences"`
}

type DefinitionRepo interface {
	GetAll(ctx context.Context, studySetID int64) ([]*Definition, error)
	Insert(ctx context.Context, studySetID int64, insertData *InsertDefinitionData) error
}

type DefinitionUseCase interface {
	Insert(ctx context.Context, userID string, studySetID int64, insertData *InsertDefinitionData) error
}

type StudySet struct {
	Id                 int64        `json:"id"`
	AuthorId           string       `json:"authorId"`
	Name               string       `json:"name"`
	Description        string       `json:"description"`
	PhraseLanguage     string       `json:"phraseLanguage"`
	DefinitionLanguage string       `json:"definitionLanguage"`
	Definitions        []Definition `json:"definitions"`
}

type StudySetSummary struct {
	Id                 int64  `json:"id"`
	AuthorId           string `json:"authorId"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	PhraseLanguage     string `json:"phraseLanguage"`
	DefinitionLanguage string `json:"definitionLanguage"`
}

type InsertStudySetData struct {
	AuthorId           string `json:"-" validate:"required"`
	Name               string `json:"name" validate:"required,ascii,max=128"`
	Description        string `json:"description" validate:"required,ascii,max=512"`
	PhraseLanguage     string `json:"phraseLanguage" validate:"required"`
	DefinitionLanguage string `json:"definitionLanguage" validate:"required"`
}

type UpdateStudySetData struct {
	Name               string `json:"name" validate:"required,ascii,max=128"`
	Description        string `json:"description" validate:"required,ascii,max=512"`
	PhraseLanguage     string `json:"phraseLanguage" validate:"required"`
	DefinitionLanguage string `json:"definitionLanguage" validate:"required"`
}

type StudySetRepo interface {
	GetAllSummary(ctx context.Context) ([]*StudySetSummary, error)
	GetById(ctx context.Context, studySetID int64) (*StudySet, error)
	Insert(ctx context.Context, insertData *InsertStudySetData) (*StudySet, error)
	Update(ctx context.Context, studySetID int64, updateData *UpdateStudySetData) error
	Delete(ctx context.Context, studySetID int64) error
}

type StudySetUseCase interface {
	GetAllSummary(ctx context.Context) ([]*StudySetSummary, error)
	GetById(ctx context.Context, studySetID int64) (*StudySet, error)
	Create(ctx context.Context, createData *InsertStudySetData) (*StudySet, error)
	Update(ctx context.Context, studySetID int64, userID string, updateData *UpdateStudySetData) error
	Delete(ctx context.Context, studySetID int64, userID string) error
}
