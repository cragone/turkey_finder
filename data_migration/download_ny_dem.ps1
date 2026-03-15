$dir = "C:\Users\Charlie\Desktop\websites\secrets-manager\turkey_finder\ny_dem_tiles"
New-Item -ItemType Directory -Force -Path $dir | Out-Null

$minLon = -79.90
$maxLon = -71.80
$minLat = 40.45
$maxLat = 45.10

$step = 0.40

for ($x = $minLon; $x -lt $maxLon; $x += $step) {
  for ($y = $minLat; $y -lt $maxLat; $y += $step) {
    $x2 = [math]::Min($x + $step, $maxLon)
    $y2 = [math]::Min($y + $step, $maxLat)

    $name = "dem_{0}_{1}.tif" -f (($x.ToString("0.00")).Replace("-","m").Replace(".","p")), (($y.ToString("0.00")).Replace("-","m").Replace(".","p"))
    $out  = Join-Path $dir $name

    $url = "https://elevation.nationalmap.gov/arcgis/rest/services/3DEPElevation/ImageServer/exportImage?bbox=$x,$y,$x2,$y2&bboxSR=4326&imageSR=4326&format=tiff&pixelType=F32&interpolation=RSP_BilinearInterpolation&f=image"

    Write-Host "Downloading $name"
    Invoke-WebRequest -Uri $url -OutFile $out
  }
}
