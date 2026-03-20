-- Spatial index on public land parcels
CREATE INDEX IF NOT EXISTS new_york_geom_idx
    ON public.new_york
    USING GIST (geom);

-- Raster convex-hull indexes (required for ST_Intersects on rasters)
CREATE INDEX IF NOT EXISTS dem_tiles_rast_idx
    ON gis.dem_tiles
    USING GIST (ST_ConvexHull(rast));

CREATE INDEX IF NOT EXISTS land_cover_rast_idx
    ON gis.land_cover
    USING GIST (ST_ConvexHull(rast));

-- Vector layer spatial indexes
CREATE INDEX IF NOT EXISTS wetlands_geom_idx
    ON gis.wetlands
    USING GIST (geom);

CREATE INDEX IF NOT EXISTS streams_geom_idx
    ON gis.streams
    USING GIST (geom);

CREATE INDEX IF NOT EXISTS nat_comm_geom_idx
    ON gis.natural_communities
    USING GIST (geom);

CREATE INDEX IF NOT EXISTS sightings_geom_idx
    ON gis.turkey_sightings
    USING GIST (geom);

CREATE INDEX IF NOT EXISTS sightings_obs_dt_idx
    ON gis.turkey_sightings (obs_dt);
