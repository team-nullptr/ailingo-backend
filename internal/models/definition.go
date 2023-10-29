package models

// Definition represents internal definition structure
type Definition struct {
	Phrase    string   `json:"phrase"`
	Meaning   string   `json:"meaning"`
	Sentences []string `json:"sentences"`
}
