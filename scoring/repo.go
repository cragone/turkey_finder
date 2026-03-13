package scoring

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type State struct {
	Land []Parcel
}

type Parcel struct {
	ID    int64
	Name  string
	Acres float64
	Score float64
}

func ListParcels(ctx context.Context, db *pgxpool.Pool) (*State, error) {
	query := `
		SELECT
			objectid AS id,
			COALESCE(unit_nm, loc_nm, 'Unnamed Parcel') AS name,
			COALESCE(gis_acres, 0) AS acres
		FROM public.new_york
		WHERE geom IS NOT NULL
		AND gis_acres > 0;
	`

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	parcels := []Parcel{}

	for rows.Next() {
		var p Parcel

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Acres,
		)
		if err != nil {
			return nil, err
		}

		parcels = append(parcels, p)
	}

	return &State{Land: parcels}, rows.Err()
}

func (state *State) ReturnParcelsWithScore() {
	for i := range state.Land {
		state.Land[i].Score = AcreageScore(state.Land[i].Acres)
	}
}
