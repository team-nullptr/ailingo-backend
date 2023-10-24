package studyset

import "ailingo/internal/models"

// StudySet represents study set structure stored in the database.
type StudySet struct {
	Id                 int64         `json:"id"`
	AuthorId           string        `json:"authorId"`
	Name               string        `json:"name"`
	Description        string        `json:"description"`
	PhraseLanguage     string        `json:"phraseLanguage"`
	DefinitionLanguage string        `json:"definitionLanguage"`
	Definitions        []models.Word `json:"definitions"`
}

// studySetCreateData represents data required to create a new study set.
// TODO: Possible values for languages are `PL` and `EN`. Can we validate that in a neat way?
type studySetCreateData struct {
	AuthorId           string        `json:"-" validate:"required"`
	Name               string        `json:"name" validate:"required,ascii,max=128"`
	Description        string        `json:"description" validate:"required,ascii,max=512"`
	PhraseLanguage     string        `json:"phraseLanguage" validate:"required,alpha,max=2"`
	DefinitionLanguage string        `json:"definitionLanguage" validate:"required,alpha,max=2"`
	Definitions        []models.Word `json:"definitions"`
}
