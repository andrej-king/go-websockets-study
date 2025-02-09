package matches

type Match struct {
	Name         string  `json:"name"`
	FirstWinOdd  float64 `json:"firstWinOdd"`
	SecondWinOdd float64 `json:"secondWinOdd"`
	DrawOdd      float64 `json:"drawOdd"`
}
