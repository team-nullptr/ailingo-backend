package chat

import (
	"ailingo/apiutil"
	"ailingo/models"
	"encoding/json"
	"net/http"
)

type Controller struct {
	sg *SentenceGenerator
}

func NewController(sg *SentenceGenerator) *Controller {
	return &Controller{
		sg: sg,
	}
}

// GenerateSentence is an endpoint handler that generates example sentences for given word.
func (c *Controller) GenerateSentence(w http.ResponseWriter, r *http.Request) {
	var word models.Word

	if err := json.NewDecoder(r.Body).Decode(&word); err != nil {
		apiutil.Err(w, http.StatusBadRequest, "expected word payload")
		return
	}

	result, err := c.sg.GenerateSentence(word)
	if err != nil {
		apiutil.Err(w, http.StatusInternalServerError, "failed to generate a sentence")
		return
	}

	apiutil.Json(w, http.StatusOK, result)
}
