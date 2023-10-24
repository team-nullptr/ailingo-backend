package studyset

import "ailingo/internal/models"

// StudySet represents study set structure stored in the database.
type StudySet struct {
	Id                 int64         `json:"id"`
	AuthorId           int64         `json:"authorId"`
	Name               string        `json:"name"`
	Description        string        `json:"description"`
	PhraseLanguage     string        `json:"phraseLanguage"`
	DefinitionLanguage string        `json:"definitionLanguage"`
	Definitions        []models.Word `json:"definitions"`
}

// StudySetCreate represents data required to create a new study set.
type StudySetCreate struct {
	AuthorId    int64  `json:"authorId" validate:"required"`
	Name        string `json:"name" validate:"required,ascii,max=128"`
	Description string `json:"description" validate:"required,ascii,max=512"`

	// TODO: Is there some neat way to validate those?
	PhraseLanguage     string `json:"phraseLanguage" validate:"required,alpha,max=2"`
	DefinitionLanguage string `json:"definitionLanguage" validate:"required,alpha,max=2"`

	Definitions []models.Word `json:"definitions"`
}
