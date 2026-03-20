package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const (
	eBirdBaseURL  = "https://api.ebird.org/v2"
	wildTurkeyCode = "wituhr"
	nyRegionCode  = "US-NY"
	maxResults    = 10000
)

type eBirdObservation struct {
	LocName string  `json:"locName"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	ObsDt   string  `json:"obsDt"` // "2024-11-01 07:30"
	HowMany int     `json:"howMany"`
}

func main() {
	_ = godotenv.Load()

	apiKey := os.Getenv("EBIRD_API_KEY")
	if apiKey == "" {
		log.Fatal("EBIRD_API_KEY env var is required")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL env var is required")
	}

	observations, err := fetchSightings(apiKey)
	if err != nil {
		log.Fatalf("fetching eBird sightings: %v", err)
	}
	fmt.Printf("Fetched %d turkey sightings from eBird\n", len(observations))

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("connecting to database: %v", err)
	}
	defer pool.Close()

	inserted, err := upsertSightings(ctx, pool, observations)
	if err != nil {
		log.Fatalf("inserting sightings: %v", err)
	}
	fmt.Printf("Inserted %d new sightings into gis.turkey_sightings\n", inserted)
}

func fetchSightings(apiKey string) ([]eBirdObservation, error) {
	url := fmt.Sprintf("%s/data/obs/%s/recent/%s?maxResults=%d",
		eBirdBaseURL, nyRegionCode, wildTurkeyCode, maxResults)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-eBirdApiToken", apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("eBird API returned status %d", resp.StatusCode)
	}

	var obs []eBirdObservation
	if err := json.NewDecoder(resp.Body).Decode(&obs); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return obs, nil
}

func upsertSightings(ctx context.Context, pool *pgxpool.Pool, obs []eBirdObservation) (int, error) {
	const query = `
		INSERT INTO gis.turkey_sightings (loc_name, lat, lng, obs_dt, how_many)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (loc_name, obs_dt) DO NOTHING
	`

	inserted := 0
	for _, o := range obs {
		// ObsDt format is "2024-11-01 07:30" or "2024-11-01"
		var obsDate *time.Time
		for _, layout := range []string{"2006-01-02 15:04", "2006-01-02"} {
			if t, err := time.Parse(layout, o.ObsDt); err == nil {
				obsDate = &t
				break
			}
		}

		tag, err := pool.Exec(ctx, query, o.LocName, o.Lat, o.Lng, obsDate, o.HowMany)
		if err != nil {
			return inserted, fmt.Errorf("inserting sighting %q: %w", o.LocName, err)
		}
		inserted += int(tag.RowsAffected())
	}
	return inserted, nil
}
