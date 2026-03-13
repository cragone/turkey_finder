package scoring

//  < 10 acres      -> 0
// 10–24 acres     -> 10
// 25–49 acres     -> 20
// 50–99 acres     -> 35
// 100–249 acres   -> 50
// 250–499 acres   -> 65
// 500–999 acres   -> 80
// 1000–4999 acres -> 90
// 5000+ acres     -> 100

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
