package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/skye-lopez/go-index.prod/idx"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Issue loading .env file --\n%s", err)
	}

	idx.FetchIdx()
}
