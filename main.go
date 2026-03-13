package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found, using system env")
	}

	log.Println("env loaded")
}
		log.Fatal(err)
	}
}
