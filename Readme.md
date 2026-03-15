this is a turkey locating app

each iteration we will build weights based on factors to score where turkeys could possibly be.

```
Invoke-WebRequest `
  -Uri "https://elevation.nationalmap.gov/arcgis/rest/services/3DEPElevation/ImageServer/exportImage?bbox=-74.30,42.60,-73.60,43.10&bboxSR=4326&imageSR=4326&format=tiff&pixelType=F32&interpolation=RSP_BilinearInterpolation&f=image" `
  -OutFile "C:\Users\Charlie\Desktop\websites\secrets-manager\turkey_finder\ny_test_dem.tif"
  ```
  
  getting the .tif for the elevation replace with correct pathways for local
