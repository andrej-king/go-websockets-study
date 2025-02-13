package main

import (
	"encoding/json"
	"go_websocket/update-page/config"
	"go_websocket/update-page/internal/handlers/api/matches"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	app := config.App{
		Port:               8080,
		IsDebug:            true,
		UpdateLiveInterval: 3 * time.Second,
		MaxOddValue:        50,
	}

	// init api matches List
	matchesApi := matches.New(&app)

	// run auto update matches odds
	go matchesApi.Run()

	// handle routes
	handleConnections(&app, matchesApi)

	log.Println("Server is started at port", app.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(app.Port), nil))
}

func handleConnections(app *config.App, matchesList *matches.List) {
	// Serve static files from the "static" directory
	http.Handle("/", http.FileServer(http.Dir("./ui")))

	// TODO  WebSocket endpoint
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(app)
	})

	// live matches endpoint (listen only get method)
	http.HandleFunc("GET /api/matches/live", func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		matchesList.Lock()
		defer matchesList.Unlock()

		json.NewEncoder(w).Encode(matchesList.Live)
	})
}
