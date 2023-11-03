package domain

import (
	"context"
)

// Author represents user information attached to study set.
type Author struct {
	Id       string `json:"id"`
	ImageURL string `json:"imageUrl"`
	Username string `json:"username"`
}

// StudySetWithAuthor represents final form of study set information.
type StudySetWithAuthor struct {
	Id                 int64  `json:"id"`
	Author             Author `json:"author"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	PhraseLanguage     string `json:"phraseLanguage"`
	DefinitionLanguage string `json:"definitionLanguage"`
}

// StudySet represents data stored in study set table.
type StudySet struct {
	Id                 int64
	Name               string
	Description        string
	PhraseLanguage     string
	DefinitionLanguage string
}

type InsertStudySetData struct {
	AuthorId           string `json:"-" validate:"required"`
	Name               string `json:"name" validate:"required,max=128"`
	Description        string `json:"description" validate:"required,max=512"`
	PhraseLanguage     string `json:"phraseLanguage" validate:"required,max=16"`
	DefinitionLanguage string `json:"definitionLanguage" validate:"required,max=16"`
}

type UpdateStudySetData struct {
	Name               string `json:"name" validate:"required,max=128"`
	Description        string `json:"description" validate:"required,max=512"`
	PhraseLanguage     string `json:"phraseLanguage" validate:"required,max=16"`
	DefinitionLanguage string `json:"definitionLanguage" validate:"required,max=16"`
}

// StudySetRepo describes methods required by StudySetRepo implementation.
type StudySetRepo interface {
	GetAll(ctx context.Context) ([]*StudySetWithAuthor, error)
	GetById(ctx context.Context, studySetID int64) (*StudySetWithAuthor, error)
	GetCreatedBy(ctx context.Context, userID string) ([]*StudySet, error)
	GetStarredBy(ctx context.Context, userID string) ([]*StudySetWithAuthor, error)
	Insert(ctx context.Context, insertData *InsertStudySetData) (int64, error)
	Update(ctx context.Context, studySetID int64, updateData *UpdateStudySetData) error
	Delete(ctx context.Context, studySetID int64) error
	Exists(ctx context.Context, studySetID int64) (bool, error)
}

// StudySetUseCase describes methods required by StudySetUseCase implementation.
type StudySetUseCase interface {
	GetAll(ctx context.Context) ([]*StudySetWithAuthor, error)
	GetById(ctx context.Context, studySetID int64) (*StudySetWithAuthor, error)
	Create(ctx context.Context, createData *InsertStudySetData) (int64, error)
	Update(ctx context.Context, userID string, studySetID int64, updateData *UpdateStudySetData) error
	Delete(ctx context.Context, userID string, studySetID int64) error
}
