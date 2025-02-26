package matches

import (
	"encoding/json"
	"go_websocket/update-page/internal/handlers/ws"
	"log"
	"time"
)

func (list *List) init(teams []string) {
	r := list.Rand
	ir := intRange{10000, 100000}
	fr := floatRange{1, list.App.MaxOddValue}

	// set init match values
	for _, v := range teams {
		match := Match{
			FirstWinOdd:  fr.NextRandomWith2Decimal(r),
			SecondWinOdd: fr.NextRandomWith2Decimal(r),
			DrawOdd:      fr.NextRandomWith2Decimal(r),
			Name:         v,
		}

		list.upsertLiveMatch(ir.NextRandom(r), match)
	}
}

type LiveOddsPayload struct {
	Name string `json:"name"`
}

// updateByInterval partial update odds by interval
func (list *List) updateLiveByInterval() {
	r := list.Rand
	ir := intRange{-10, 10}                   // for partial update random odds (no all)
	fr := floatRange{1, list.App.MaxOddValue} // for update new odd values
	var oldMatchData Match
	var isMatchUpdated bool

	for {
		time.Sleep(list.App.UpdateLiveInterval)

		// update random odds in live matches
		for k, v := range list.Live {
			// set default values
			oldMatchData = v
			isMatchUpdated = false

			// update first win odd
			if ir.IsOddUpdateNeeded(r) {
				isMatchUpdated = true

				// override match value in range copy
				v.FirstWinOdd = fr.NextRandomWith2Decimal(r)

				// override match in original map
				list.Live[k] = v
			}

			// update second win odd
			if ir.IsOddUpdateNeeded(r) {
				isMatchUpdated = true

				// override match value in range copy
				v.SecondWinOdd = fr.NextRandomWith2Decimal(r)

				// override match in original map
				list.Live[k] = v
			}

			// update draw odd
			if ir.IsOddUpdateNeeded(r) {
				isMatchUpdated = true

				// override match value in range copy
				v.DrawOdd = fr.NextRandomWith2Decimal(r)

				// override match in original map
				list.Live[k] = v
			}

			// send updated odds as json
			if isMatchUpdated {
				jsonString, _ := json.Marshal(v)
				list.WebsocketManager.SubscribersDataHandler(ws.Event{Type: ws.EventLiveOdds, Payload: jsonString})
			}

			// show changes in console
			if isMatchUpdated && list.App.IsDebug {
				log.Printf("[%d] match updated\n", k)
				log.Println("Old:", oldMatchData)
				log.Println("New:", v)
				log.Println("-------------------------------------")
			}
		}
	}
}
