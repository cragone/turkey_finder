package main

import (
	"char/db"
	"char/scoring"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panic("no .env file found, using system env")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Panic("no database url found")
	}

	appConn, err := db.New(context.Background(), dsn)
	if err != nil {
		panic(err)
	}

	newYork, err := scoring.ListParcels(context.Background(), appConn.Pool)
	if err != nil {
		log.Fatal(err)
	}

	newYork.ReturnParcelsWithScore()

	for _, p := range newYork.Land {
		fmt.Printf("Score: %f | %s | %.1f acres\n", p.Score, p.Name, p.Acres)
	}
}
