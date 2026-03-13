package main

import (
	"char/db"
	"context"
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

	_ = appConn
}
