package models

// StudySet represents study set structure exposed to the end user.
type StudySet struct {
	Id                 int64  `json:"id"`
	AuthorId           string `json:"authorId"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	PhraseLanguage     string `json:"phraseLanguage"`
	DefinitionLanguage string `json:"definitionLanguage"`
	Definitions        []Word `json:"definitions"`
}

// StudySetSummary represents study set structure exposed to the end user that can be used for rendering in lists.
type StudySetSummary struct {
	Id                 int64  `json:"id"`
	AuthorId           string `json:"authorId"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	PhraseLanguage     string `json:"phraseLanguage"`
	DefinitionLanguage string `json:"definitionLanguage"`
}
