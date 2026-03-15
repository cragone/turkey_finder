# Turkey Finder

_This is used to find and grade public land parcels for the presence of turkey._



## Getting Elevation Data
### getting a .tif for the elevation replace with correct pathways for local
```
Invoke-WebRequest `
  -Uri "https://elevation.nationalmap.gov/arcgis/rest/services/3DEPElevation/ImageServer/exportImage?bbox=-74.30,42.60,-73.60,43.10&bboxSR=4326&imageSR=4326&format=tiff&pixelType=F32&interpolation=RSP_BilinearInterpolation&f=image" `
  -OutFile "C:\Users\Charlie\Desktop\websites\secrets-manager\turkey_finder\ny_test_dem.tif"
  ```
  
- see data_migration for script of whole state will take a while for full upload.

```
Get-ChildItem "C:\Users\Charlie\Desktop\websites\secrets-manager\turkey_finder\ny_dem_tiles\*.tif" |
ForEach-Object {
  & raster2pgsql -a -I -C -M -F -t 256x256 $_.FullName public.ny_dem
} | & "C:\Program Files\PostgreSQL\18\bin\psql.exe" -d turkey_finder -U postgres
```
