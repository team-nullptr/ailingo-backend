package deepl

// DeeplTranslationRequest represents a deepl translation api request payload.
type DeeplTranslationRequest struct {
	Text       []string `json:"text"`
	TargetLang string   `json:"target_lang"`
}

// Translation represents a translation entry from DeeplTranslationResult.
type Translation struct {
	DetectedSourceLanguage string `json:"detected_source_language"`
	Text                   string `json:"text"`
}

// DeeplTranslationResult represents a response from translation api.
type DeeplTranslationResult struct {
	Translations []Translation `json:"translations"`
}
