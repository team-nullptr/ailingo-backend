package chat

import (
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

// GenerateSentence is an endpoint handler that generates example sentences for given definition.
func (c *Controller) GenerateSentence(w http.ResponseWriter, r *http.Request) {
	var definition models.Definition

	if err := json.NewDecoder(r.Body).Decode(&definition); err != nil {
		// TODO: handle error
		panic(err)
	}

	sentence, err := c.sg.GenerateSentence(definition)
	if err != nil {
		// TODO: handle error
		panic(err)
	}

	w.Write([]byte(sentence))
}
