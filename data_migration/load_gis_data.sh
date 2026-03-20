#!/usr/bin/env bash
# load_gis_data.sh
#
# Instructions for loading all NY GIS layers into PostGIS.
# Prerequisites: psql, shp2pgsql, raster2pgsql, gdal (gdal_translate, gdalwarp)
#
# Set these before running:
#   export PGDATABASE=turkey_finder
#   export PGUSER=postgres
#   export PGHOST=localhost
#   export PGPORT=5432
#
# Run the schema first:
#   psql -f database_schema.sql
#   psql -f indexes.sql
#
# ---------------------------------------------------------------------------

set -euo pipefail

PSQL="psql -d ${PGDATABASE:-turkey_finder} -U ${PGUSER:-postgres} -h ${PGHOST:-localhost} -p ${PGPORT:-5432}"
SHP2PGSQL="shp2pgsql"
RASTER2PGSQL="raster2pgsql"

# ---------------------------------------------------------------------------
# 1. NWI WETLANDS
#    Source: https://www.fws.gov/program/national-wetlands-inventory/wetlands-mapper
#    Select "New York" -> Download Shapefile
#    File: NY_Wetlands.zip  ->  NY_Wetlands.shp
# ---------------------------------------------------------------------------
echo "=== Loading NWI Wetlands ==="
# Unzip if needed: unzip NY_Wetlands.zip
# Reproject to WGS84 (4326) if source is not already in 4326:
#   ogr2ogr -f "ESRI Shapefile" -t_srs EPSG:4326 NY_Wetlands_4326.shp NY_Wetlands.shp
$SHP2PGSQL -s 4326 -a -I NY_Wetlands.shp gis.wetlands | $PSQL
echo "Wetlands loaded."

# ---------------------------------------------------------------------------
# 2. NHD STREAMS
#    Source: https://www.usgs.gov/national-hydrography/access-national-hydrography-products
#    Select "New York" -> NHD Best Resolution -> NHDFlowline
#    File: NHD_H_New_York_State_Shape.zip  ->  Shape/NHDFlowline.shp
# ---------------------------------------------------------------------------
echo "=== Loading NHD Streams ==="
# Reproject to WGS84 if needed:
#   ogr2ogr -f "ESRI Shapefile" -t_srs EPSG:4326 NHDFlowline_4326.shp Shape/NHDFlowline.shp
$SHP2PGSQL -s 4326 -a -I NHDFlowline.shp gis.streams | $PSQL
echo "Streams loaded."

# ---------------------------------------------------------------------------
# 3. NLCD LAND COVER
#    Source: https://www.mrlc.gov/data
#    Product: NLCD 2021 Land Cover (CONUS) -> download GeoTIFF
#    File: nlcd_2021_land_cover_l48_20230630.img (or .tif)
#
#    Forest NLCD class codes:
#      41 = Deciduous Forest
#      42 = Evergreen Forest
#      43 = Mixed Forest
#
#    Step 1: Clip to NY bounding box (-79.9 to -71.8, 40.4 to 45.2)
#    Step 2: Load clipped raster into PostGIS
# ---------------------------------------------------------------------------
echo "=== Loading NLCD Land Cover ==="
# Clip to NY extent (adjust input filename as needed):
gdalwarp \
  -te -79.90 40.40 -71.80 45.20 \
  -t_srs EPSG:4326 \
  -of GTiff \
  nlcd_2021_land_cover_l48_20230630.img \
  nlcd_ny_clipped.tif

$RASTER2PGSQL -a -I -C -M -F -t 256x256 nlcd_ny_clipped.tif gis.land_cover | $PSQL
echo "NLCD land cover loaded."

# ---------------------------------------------------------------------------
# 4. DEM ELEVATION TILES
#    See data_migration/download_ny_dem.ps1 for tile download.
#    Assumes tiles are in ./ny_dem_tiles/*.tif
# ---------------------------------------------------------------------------
echo "=== Loading DEM Elevation Tiles ==="
for f in ny_dem_tiles/*.tif; do
  $RASTER2PGSQL -a -I -C -M -F -t 256x256 "$f" gis.dem_tiles | $PSQL
done
echo "DEM tiles loaded."

# ---------------------------------------------------------------------------
# 5. NY NATURAL HERITAGE (optional)
#    Source: https://guides.nynhp.org/downloads/
#    File: ny_natural_communities.shp (if available)
# ---------------------------------------------------------------------------
# echo "=== Loading NY Natural Heritage ==="
# $SHP2PGSQL -s 4326 -a ny_natural_communities.shp gis.natural_communities | $PSQL
# echo "Natural communities loaded."

echo ""
echo "All GIS data loaded. Run ANALYZE to update statistics:"
echo "  psql -c 'ANALYZE gis.dem_tiles; ANALYZE gis.land_cover; ANALYZE gis.wetlands; ANALYZE gis.streams;'"
