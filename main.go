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
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system env")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL env var is required")
	}

	appConn, err := db.New(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer appConn.Close()

	newYork, err := scoring.ListParcels(context.Background(), appConn.Pool)
	if err != nil {
		log.Fatal(err)
	}

	newYork.ReturnParcelsWithScore()

	fmt.Printf("%-6s  %-50s  %8s  %8s  %7s  %9s  %8s  %5s\n",
		"Score", "Name", "Acres", "ElevσM", "Forest%", "Water(m)", "Sightings", "")
	fmt.Println("------  --------------------------------------------------  --------  --------  -------  ---------  ---------")

	for _, p := range newYork.Land {
		fmt.Printf("%6.1f  %-50s  %8.1f  %8.1f  %6.1f%%  %9.0f  %9d\n",
			p.Score,
			p.Name,
			p.Acres,
			p.ElevStdDev,
			p.ForestPct*100,
			p.WaterDistMeters,
			p.SightingCount,
		)
	}
}
