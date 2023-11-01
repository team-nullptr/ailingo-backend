package domain

import (
	"context"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

// Author represents user information attached to study set.
type Author struct {
	ProfileImageURL string `json:"profileImageUrl"`
	Username        string `json:"username"`
}

// AuthorFromClerkUser creates a new Author for the given clerk.User.
func AuthorFromClerkUser(user *clerk.User) *Author {
	var username string

	if user.Username != nil {
		username = *user.Username
	}

	return &Author{
		ProfileImageURL: user.ProfileImageURL,
		Username:        username,
	}
}

// StudySet represents final form of study set information.
type StudySet struct {
	Id                 int64   `json:"id"`
	Author             *Author `json:"author"`
	Name               string  `json:"name"`
	Description        string  `json:"description"`
	PhraseLanguage     string  `json:"phraseLanguage"`
	DefinitionLanguage string  `json:"definitionLanguage"`
}

// StudySetRow represents data stored in study set table.
type StudySetRow struct {
	Id                 int64
	AuthorId           string
	Name               string
	Description        string
	PhraseLanguage     string
	DefinitionLanguage string
}

// Populate populates study set information with that is not stored in study set table.
func (r *StudySetRow) Populate(author *Author) *StudySet {
	return &StudySet{
		Id:                 r.Id,
		Author:             author,
		Name:               r.Name,
		Description:        r.Description,
		PhraseLanguage:     r.PhraseLanguage,
		DefinitionLanguage: r.DefinitionLanguage,
	}
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
	GetAll(ctx context.Context) ([]*StudySetRow, error)
	GetById(ctx context.Context, studySetID int64) (*StudySetRow, error)
	Insert(ctx context.Context, insertData *InsertStudySetData) (*StudySetRow, error)
	Update(ctx context.Context, studySetID int64, updateData *UpdateStudySetData) error
	Delete(ctx context.Context, studySetID int64) error
	Exists(ctx context.Context, studySetID int64) (bool, error)
}

// StudySetUseCase describes methods required by StudySetUseCase implementation.
type StudySetUseCase interface {
	GetAll(ctx context.Context) ([]*StudySet, error)
	GetById(ctx context.Context, studySetID int64) (*StudySet, error)
	Create(ctx context.Context, createData *InsertStudySetData) (*StudySet, error)
	Update(ctx context.Context, studySetID int64, userID string, updateData *UpdateStudySetData) error
	Delete(ctx context.Context, studySetID int64, userID string) error
}
