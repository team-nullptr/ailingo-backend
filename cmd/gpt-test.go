package main

import (
	"ailingo/gpt"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("failed to load configuration from .env file\n")
	}

	service := gpt.NewExampleSentenceService(http.DefaultClient)

	s, err := service.GenerateSentence(gpt.Definition{
		Phrase:  "judicial branch",
		Meaning: "władza sądownicza",
	})
	if err != nil {
		log.Fatal("generating sentence failed: %w", err)
	}

	fmt.Println(s)
}
