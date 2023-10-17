package models

type Definition struct {
	Phrase      string `json:"phrase"`
	PhraseLang  string `json:"phraseLang"`
	Meaning     string `json:"meaning"`
	MeaningLang string `json:"meaningLang"`
}
