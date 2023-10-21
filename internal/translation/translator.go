package translation

import "context"

type Translator interface {
	Translate(ctx context.Context, text string) (string, error)
}
