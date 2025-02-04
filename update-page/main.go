package main

import (
	"encoding/json"
	"go_websocket/update-page/api"
	"log"
	"net/http"
	"strconv"
	"time"
)

var port = 8000
var updateMatchesInterval = 3 * time.Second

func main() {
	handleRoutes()

	log.Printf("Going to listen on port %d\n", port)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(port), nil))
}

func handleRoutes() {
	// init match updates
	go api.Init(updateMatchesInterval)

	http.Handle("/", http.FileServer(http.Dir("./public")))

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("{}")
	})

	http.HandleFunc("/api/matches", func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(api.MatchList) // or use json.Marshal
	})
}
