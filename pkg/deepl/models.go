package deepl

// TranslationRequest represents a deepl translation api request payload.
type TranslationRequest struct {
	Text       []string `json:"text"`
	TargetLang string   `json:"target_lang"`
}

// Translation represents a translation entry from DeeplTranslationResult.
type Translation struct {
	DetectedSourceLanguage string `json:"detected_source_language"`
	Text                   string `json:"text"`
}

// TranslationResult represents a response from translation api.
type TranslationResult struct {
	Translations []Translation `json:"translations"`
}
