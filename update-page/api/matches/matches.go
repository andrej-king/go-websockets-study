package matches

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// example match names
var teams = []string{
	"Mansfield Town - Northampton",
	"Burton Albion - Blackpool",
	"Rotherham - Shrewsbury",
	"Barcelona - Alaves",
	"Santos - Sao Paulo",
	"Indiana Pacers - Atlanta Hawks",
	"Vitoria Guimaraes - AVS Sad",
	"Real Santander - Tigres",
	"Deportiva Venados - San-Juan de Aragon",
	"Atletico Nayarit - Puerto Vallarta",
	"Vitoria BA - Real Noroeste",
	"Ecuador (U20) - Argentina (U20)",
}

var MatchList = map[int]Match{}

type Match struct {
	Name string  `json:"name"`
	Odd  float64 `json:"odd"`
}

type IntRange struct {
	min, max int
}

type FloatRange struct {
	min, max float64
}

// NextRandom random int value in range
func (ir *IntRange) NextRandom(r *rand.Rand) int {
	return r.Intn(ir.max-ir.min+1) + ir.min
}

// NextRandom random float64 value in range
func (fr *FloatRange) NextRandom(r *rand.Rand) float64 {
	return fr.min + r.Float64()*(fr.max-fr.min)
}

// NextRandomWith2Decimal random float64 value in range with 2 decimal
func (fr *FloatRange) NextRandomWith2Decimal(r *rand.Rand) float64 {
	result := fr.NextRandom(r)

	return math.Floor(result*100) / 100
}

// initMatchesData set initial match values
func initMatchesData() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ir := IntRange{10000, 100000}

	// set init match values
	for _, v := range teams {
		MatchList[ir.NextRandom(r)] = Match{Odd: 0, Name: v}
	}
}

// updateOdds partial update odds by interval
func updateMatches(updateInterval time.Duration) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ir := IntRange{-10, 10} // for partial update random odds (no all)
	fr := FloatRange{1, 50} // for update new odd values
	var oldOdd float64

	for {
		// update random odds in matches
		for k, v := range MatchList {
			// update all odds for first start
			if v.Odd <= 1 || ir.NextRandom(r) > 0 {
				oldOdd = v.Odd

				// generate a new random odd value and set as new map value
				v.Odd = fr.NextRandomWith2Decimal(r)
				MatchList[k] = v

				fmt.Printf("[%d] match with updated odds: [%.2f] old, [%.2f] new\n", k, oldOdd, v.Odd)
			}
		}

		fmt.Println()

		time.Sleep(updateInterval)
	}
}

func Init(updateInterval time.Duration) {
	initMatchesData()
	updateMatches(updateInterval)
}
