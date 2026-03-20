package scoring

// Composite score weights (must sum to 1.0)
const (
	weightAcreage  = 0.15
	weightElev     = 0.20
	weightForest   = 0.25
	weightWater    = 0.20
	weightSighting = 0.20
)

// AcreageScore scores a parcel on size alone.
//
//	 < 10 acres      ->   0
//	10–24 acres     ->  10
//	25–49 acres     ->  20
//	50–99 acres     ->  35
//	100–249 acres   ->  50
//	250–499 acres   ->  65
//	500–999 acres   ->  80
//	1000–4999 acres ->  90
//	5000+ acres     -> 100
func AcreageScore(acres float64) float64 {
	switch {
	case acres < 10:
		return 0
	case acres < 25:
		return 10
	case acres < 50:
		return 20
	case acres < 100:
		return 35
	case acres < 250:
		return 50
	case acres < 500:
		return 65
	case acres < 1000:
		return 80
	case acres < 5000:
		return 90
	default:
		return 100
	}
}

// ElevationScore scores terrain ruggedness from DEM standard deviation.
// Turkeys prefer rolling hills for roosting and escape cover.
//
//	 0– 5 m stddev ->   0 (flat, poor)
//	 5–15 m        ->  20
//	15–30 m        ->  40
//	30–50 m        ->  60
//	50–100 m       ->  80
//	100+ m         -> 100 (highly varied, excellent)
func ElevationScore(stddevMeters float64) float64 {
	switch {
	case stddevMeters < 5:
		return 0
	case stddevMeters < 15:
		return 20
	case stddevMeters < 30:
		return 40
	case stddevMeters < 50:
		return 60
	case stddevMeters < 100:
		return 80
	default:
		return 100
	}
}

// ForestScore scores woodland cover fraction from NLCD (classes 41/42/43).
// Turkeys need edge habitat — mixed open and forest — not 100% closed canopy.
//
//	0.0–0.1  ->   0 (too open)
//	0.1–0.3  ->  20
//	0.3–0.5  ->  60
//	0.5–0.7  -> 100 (ideal mixed edge)
//	0.7–0.85 ->  80
//	0.85+    ->  50 (dense forest, limited foraging)
func ForestScore(pct float64) float64 {
	switch {
	case pct < 0.10:
		return 0
	case pct < 0.30:
		return 20
	case pct < 0.50:
		return 60
	case pct < 0.70:
		return 100
	case pct < 0.85:
		return 80
	default:
		return 50
	}
}

// WaterProximityScore scores proximity to the nearest wetland or stream.
// Turkeys drink daily and forage along waterways.
//
//	>5000 m   ->   0
//	2000–5000 ->  20
//	1000–2000 ->  50
//	 500–1000 ->  75
//	 100–500  ->  90
//	   0–100  -> 100 (water within the parcel or on its edge)
func WaterProximityScore(distMeters float64) float64 {
	switch {
	case distMeters > 5000:
		return 0
	case distMeters > 2000:
		return 20
	case distMeters > 1000:
		return 50
	case distMeters > 500:
		return 75
	case distMeters > 100:
		return 90
	default:
		return 100
	}
}

// SightingDensityScore scores based on verified eBird Wild Turkey observations
// within 1000 m of the parcel.
//
//	 0 sightings ->   0
//	 1–2         ->  20
//	 3–5         ->  40
//	 6–10        ->  60
//	11–25        ->  80
//	25+          -> 100
func SightingDensityScore(count int) float64 {
	switch {
	case count == 0:
		return 0
	case count <= 2:
		return 20
	case count <= 5:
		return 40
	case count <= 10:
		return 60
	case count <= 25:
		return 80
	default:
		return 100
	}
}

// CompositeScore returns a weighted habitat-suitability grade from 0–100.
//
//	Acreage          15%
//	Elevation var    20%
//	Forest coverage  25%
//	Water proximity  20%
//	Sighting density 20%
func CompositeScore(p Parcel) float64 {
	return AcreageScore(p.Acres)*weightAcreage +
		ElevationScore(p.ElevStdDev)*weightElev +
		ForestScore(p.ForestPct)*weightForest +
		WaterProximityScore(p.WaterDistMeters)*weightWater +
		SightingDensityScore(p.SightingCount)*weightSighting
}
