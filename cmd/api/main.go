package main

import (
	"flag"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"ailingo/config"
	"ailingo/internal/app"
)

func main() {
	useDotenv := flag.Bool("dotenv", false, "configure with .env")
	flag.Parse()

	cfg, err := config.New(*useDotenv)
	if err != nil {
		log.Fatal("failed to load configuration", err)
	}

	app.Run(cfg)
}
