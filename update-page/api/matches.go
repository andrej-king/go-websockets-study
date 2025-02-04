package api

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
	Name         string  `json:"name"`
	FirstWinOdd  float64 `json:"firstWinOdd"`
	SecondWinOdd float64 `json:"secondWinOdd"`
	DrawOdd      float64 `json:"drawOdd"`
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
	fr := FloatRange{1, 50}

	// set init match values
	for _, v := range teams {
		MatchList[ir.NextRandom(r)] = Match{
			FirstWinOdd:  fr.NextRandomWith2Decimal(r),
			SecondWinOdd: fr.NextRandomWith2Decimal(r),
			DrawOdd:      fr.NextRandomWith2Decimal(r),
			Name:         v,
		}
	}
}

// updateOdds partial update odds by interval
func updateMatches(updateInterval time.Duration, isDebugMatches bool) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ir := IntRange{-10, 10} // for partial update random odds (no all)
	fr := FloatRange{1, 50} // for update new odd values
	var oldMatchData Match
	var isMatchUpdated bool

	for {
		time.Sleep(updateInterval)

		// update random odds in matches
		for k, v := range MatchList {

			// set default values
			oldMatchData = v
			isMatchUpdated = false

			// update first win odd
			if ir.NextRandom(r) > 0 {
				isMatchUpdated = true
				v.FirstWinOdd = fr.NextRandomWith2Decimal(r)
				MatchList[k] = v
			}

			// update second win odd
			if ir.NextRandom(r) > 0 {
				isMatchUpdated = true
				v.SecondWinOdd = fr.NextRandomWith2Decimal(r)
				MatchList[k] = v
			}

			// update draw odd
			if ir.NextRandom(r) > 0 {
				isMatchUpdated = true
				v.DrawOdd = fr.NextRandomWith2Decimal(r)
				MatchList[k] = v
			}

			// show changes in console
			if isMatchUpdated && isDebugMatches {
				fmt.Printf("[%d] match updated\n", k)
				fmt.Println("Old:", oldMatchData)
				fmt.Println("New:", v)
				fmt.Println()
			}
		}
	}
}

func Init(updateInterval time.Duration, isDebugMatches bool) {
	initMatchesData()
	updateMatches(updateInterval, isDebugMatches)
}
