package scoring

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type State struct {
	Land []Parcel
}

type Parcel struct {
	ID              int64
	Name            string
	Acres           float64
	ElevStdDev      float64 // elevation std-dev in meters within parcel footprint
	ForestPct       float64 // fraction 0.0–1.0 of parcel covered by NLCD forest classes 41/42/43
	WaterDistMeters float64 // distance in meters to nearest wetland or stream
	SightingCount   int     // eBird Wild Turkey sightings within 1000 m of parcel
	Score           float64 // composite habitat score
}

// ListParcels fetches all scored parcels from PostGIS.
// Elevation, forest coverage, water proximity, and sighting density are
// computed via spatial joins against the gis.* layers.
// Parcels without raster/vector coverage fall back to 0 / 9999.
func ListParcels(ctx context.Context, db *pgxpool.Pool) (*State, error) {
	query := `
		WITH parcel AS (
			SELECT
				objectid                                        AS id,
				COALESCE(unit_nm, loc_nm, 'Unnamed Parcel')    AS name,
				COALESCE(gis_acres, 0)                         AS acres,
				geom
			FROM public.new_york
			WHERE geom IS NOT NULL
			  AND gis_acres > 0
		),

		-- Elevation standard deviation from DEM tiles clipped to each parcel.
		-- Uses the first intersecting tile; parcels spanning multiple tiles are
		-- approximated by the tile with the largest overlap.
		elev AS (
			SELECT DISTINCT ON (p.id)
				p.id,
				COALESCE(
					(ST_SummaryStats(ST_Clip(d.rast, p.geom, TRUE))).stddev,
					0
				) AS elev_stddev
			FROM parcel p
			JOIN gis.dem_tiles d ON ST_Intersects(d.rast, p.geom)
			ORDER BY p.id, ST_Area(ST_Intersection(ST_ConvexHull(d.rast), p.geom)) DESC
		),

		-- Forest coverage fraction from NLCD land cover raster.
		-- NLCD classes 41 (Deciduous), 42 (Evergreen), 43 (Mixed Forest).
		forest AS (
			SELECT
				p.id,
				COALESCE(
					SUM(CASE WHEN vc.val IN (41, 42, 43) THEN vc.count ELSE 0 END)::float /
					NULLIF(SUM(vc.count)::float, 0),
					0
				) AS forest_pct
			FROM parcel p
			JOIN gis.land_cover lc ON ST_Intersects(lc.rast, p.geom),
			LATERAL ST_ValueCount(ST_Clip(lc.rast, p.geom, TRUE)) AS vc(val, count)
			GROUP BY p.id
		),

		-- Distance in meters to the nearest wetland polygon or stream line.
		-- Uses geography cast so the result is in metres regardless of CRS.
		-- KNN search with <-> is index-accelerated.
		water AS (
			SELECT p.id, MIN(dist_m) AS water_dist_m
			FROM parcel p,
			LATERAL (
				(
					SELECT ST_Distance(p.geom::geography, w.geom::geography) AS dist_m
					FROM gis.wetlands w
					ORDER BY p.geom <-> w.geom
					LIMIT 1
				)
				UNION ALL
				(
					SELECT ST_Distance(p.geom::geography, s.geom::geography) AS dist_m
					FROM gis.streams s
					ORDER BY p.geom <-> s.geom
					LIMIT 1
				)
			) nearest
			GROUP BY p.id
		),

		-- Count of eBird Wild Turkey sightings within 1000 m of the parcel boundary.
		sightings AS (
			SELECT
				p.id,
				COUNT(ts.id)::int AS sighting_count
			FROM parcel p
			LEFT JOIN gis.turkey_sightings ts
				ON ST_DWithin(p.geom::geography, ts.geom::geography, 1000)
			GROUP BY p.id
		)

		SELECT
			p.id,
			p.name,
			p.acres,
			COALESCE(e.elev_stddev,     0)    AS elev_stddev,
			COALESCE(f.forest_pct,      0)    AS forest_pct,
			COALESCE(w.water_dist_m, 9999)    AS water_dist_m,
			COALESCE(si.sighting_count,  0)   AS sighting_count
		FROM parcel p
		LEFT JOIN elev      e  ON e.id  = p.id
		LEFT JOIN forest    f  ON f.id  = p.id
		LEFT JOIN water     w  ON w.id  = p.id
		LEFT JOIN sightings si ON si.id = p.id;
	`

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parcels []Parcel

	for rows.Next() {
		var p Parcel
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Acres,
			&p.ElevStdDev,
			&p.ForestPct,
			&p.WaterDistMeters,
			&p.SightingCount,
		); err != nil {
			return nil, err
		}
		parcels = append(parcels, p)
	}

	return &State{Land: parcels}, rows.Err()
}

func (state *State) ReturnParcelsWithScore() {
	for i := range state.Land {
		state.Land[i].Score = CompositeScore(state.Land[i])
	}
}
