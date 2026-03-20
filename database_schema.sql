-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS postgis_raster;

-- Create gis schema for raster/vector layers
CREATE SCHEMA IF NOT EXISTS gis;

-- ---------------------------------------------------------------------------
-- public.new_york
-- Public land parcels / PAD (Protected Areas Database) polygons for NY.
-- Load via shp2pgsql from PAD-US or NY Open Data parcel shapefile.
-- ---------------------------------------------------------------------------
-- (Table is loaded externally; ensure it has at minimum these columns)
-- objectid   SERIAL PRIMARY KEY
-- unit_nm    TEXT          -- unit name
-- loc_nm     TEXT          -- location name
-- gis_acres  DOUBLE PRECISION
-- geom       GEOMETRY(MultiPolygon, 4326)

-- ---------------------------------------------------------------------------
-- gis.dem_tiles
-- Digital Elevation Model raster tiles loaded from USGS 3DEP via raster2pgsql.
-- See data_migration/download_ny_dem.ps1 for download instructions.
-- Load: raster2pgsql -a -I -C -M -F -t 256x256 *.tif gis.dem_tiles | psql ...
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS gis.dem_tiles (
    rid      SERIAL PRIMARY KEY,
    rast     RASTER,
    filename TEXT
);

-- ---------------------------------------------------------------------------
-- gis.land_cover
-- NLCD (National Land Cover Database) raster for NY.
-- See data_migration/load_gis_data.sh for download and load instructions.
-- Forest classes: 41 (Deciduous), 42 (Evergreen), 43 (Mixed Forest)
-- Load: raster2pgsql -a -I -C -M -F -t 256x256 nlcd_ny.tif gis.land_cover | psql ...
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS gis.land_cover (
    rid      SERIAL PRIMARY KEY,
    rast     RASTER,
    filename TEXT
);

-- ---------------------------------------------------------------------------
-- gis.wetlands
-- NWI (National Wetlands Inventory) polygons for NY from USFWS.
-- See data_migration/load_gis_data.sh for download and load instructions.
-- Load: shp2pgsql -s 4326 -a NY_Wetlands.shp gis.wetlands | psql ...
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS gis.wetlands (
    gid          SERIAL PRIMARY KEY,
    wetland_type TEXT,
    geom         GEOMETRY(MultiPolygon, 4326)
);

-- ---------------------------------------------------------------------------
-- gis.streams
-- NHD (National Hydrography Dataset) flowlines for NY from USGS.
-- See data_migration/load_gis_data.sh for download and load instructions.
-- Load: shp2pgsql -s 4326 -a NHDFlowline.shp gis.streams | psql ...
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS gis.streams (
    gid       SERIAL PRIMARY KEY,
    gnis_name TEXT,
    geom      GEOMETRY(MultiLineString, 4326)
);

-- ---------------------------------------------------------------------------
-- gis.natural_communities
-- NY Natural Heritage Program community polygons.
-- Load: shp2pgsql -s 4326 -a ny_natural_communities.shp gis.natural_communities | psql ...
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS gis.natural_communities (
    gid       SERIAL PRIMARY KEY,
    comm_name TEXT,
    geom      GEOMETRY(MultiPolygon, 4326)
);

-- ---------------------------------------------------------------------------
-- gis.turkey_sightings
-- Wild Turkey observation points pulled from eBird API (species: wituhr).
-- Populated by data_migration/fetch_ebird.go.
-- geom is auto-generated from lat/lng.
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS gis.turkey_sightings (
    id        SERIAL PRIMARY KEY,
    loc_name  TEXT,
    lat       DOUBLE PRECISION NOT NULL,
    lng       DOUBLE PRECISION NOT NULL,
    obs_dt    DATE,
    how_many  INT,
    geom      GEOMETRY(Point, 4326)
                  GENERATED ALWAYS AS (ST_SetSRID(ST_MakePoint(lng, lat), 4326)) STORED,
    UNIQUE (loc_name, obs_dt)
);
