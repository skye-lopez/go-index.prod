package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/skye-lopez/go-index.prod/cmd"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Issue loading .env file \n%s", err)
	}

	cmd.Execute()
}
