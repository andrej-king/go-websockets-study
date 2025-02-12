package matches

import (
	"go_websocket/update-page/config"
	"math"
	"math/rand"
	"sync"
	"time"
)

// List keep map with updatable matches
type List struct {
	sync.RWMutex
	App  *config.App
	Rand *rand.Rand

	Live  map[int]Match
	Teams []string
}

// upsertLiveMatch update or add Match in Live List
func (list *List) upsertLiveMatch(id int, match Match) {
	list.Lock()
	defer list.Unlock()

	list.Live[id] = match
}

// New initialize matches List with default values
func New(app *config.App) *List {
	list := List{
		App:  app,
		Rand: rand.New(rand.NewSource(time.Now().UnixNano())),

		Live:  make(map[int]Match),
		Teams: teams,
	}

	return &list
}

// Run update matches odds by interval
func (list *List) Run() {
	list.init(list.Teams)
	list.updateLiveByInterval()
}

type Match struct {
	Name         string  `json:"name"`
	FirstWinOdd  float64 `json:"firstWinOdd"`
	SecondWinOdd float64 `json:"secondWinOdd"`
	DrawOdd      float64 `json:"drawOdd"`
}

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

type intRange struct {
	min, max int
}

type floatRange struct {
	min, max float64
}

// NextRandom random int value in range
func (ir *intRange) NextRandom(r *rand.Rand) int {
	return r.Intn(ir.max-ir.min+1) + ir.min
}

// IsOddUpdateNeeded helper func for random update match odd (regenerate status every call)
func (ir *intRange) IsOddUpdateNeeded(r *rand.Rand) bool {
	return ir.NextRandom(r) > 0
}

// NextRandom random float64 value in range
func (fr *floatRange) NextRandom(r *rand.Rand) float64 {
	return fr.min + r.Float64()*(fr.max-fr.min)
}

// NextRandomWith2Decimal random float64 value in range with 2 decimal
func (fr *floatRange) NextRandomWith2Decimal(r *rand.Rand) float64 {
	result := fr.NextRandom(r)

	return math.Floor(result*100) / 100
}
