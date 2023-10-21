package translation

import "context"

// DevTranslator stubs DeepL client functionality in order to save API calls during app development.
type DevTranslator struct{}

func NewDevTranslator() Translator {
	return &DevTranslator{}
}

func (dt *DevTranslator) Translate(ctx context.Context, text string) (string, error) {
	return "Development", nil
}
