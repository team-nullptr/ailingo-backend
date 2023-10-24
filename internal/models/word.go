package models

import "fmt"

// TODO: This package will probably disappear

// Word represents internal word structure
type Word struct {
	Phrase     string `json:"phrase"`
	Definition string `json:"definition"`
}

func (w *Word) ToChatPrompt() string {
	return fmt.Sprintf("english phrase: %s\npolish definition: %s", w.Phrase, w.Definition)
}
