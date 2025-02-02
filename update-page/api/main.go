package main

import (
	"encoding/json"
	"go_websocket/update-page/api/matches"
	"log"
	"net/http"
	"strconv"
	"time"
)

var port = 8000
var updateMatchesInterval = 7 * time.Second

func main() {
	go matches.Init(updateMatchesInterval)

	http.HandleFunc("/api/matches", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(matches.MatchList) // or use json.Marshal
	})

	log.Printf("Going to listen on port %d\n", port)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(port), nil))
}
