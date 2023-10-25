package models

// StudySet represents study set structure stored in the database.
type StudySet struct {
	Id                 int64  `json:"id"`
	AuthorId           string `json:"authorId"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	PhraseLanguage     string `json:"phraseLanguage"`
	DefinitionLanguage string `json:"definitionLanguage"`
	Definitions        []Word `json:"definitions"`
}
