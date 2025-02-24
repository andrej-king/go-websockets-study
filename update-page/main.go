package main

import (
	"encoding/json"
	"go_websocket/update-page/config"
	"go_websocket/update-page/internal/handlers/api/matches"
	"go_websocket/update-page/internal/handlers/ws"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	app := config.App{
		Port:               8080,
		IsDebug:            false,
		UpdateLiveInterval: 3 * time.Second,
		MaxOddValue:        50,
	}

	// init
	matchesApi := matches.New(&app)
	wsManager := ws.NewManager(&app)

	// run auto update matches odds
	go matchesApi.Run()
	//go ws.New()

	// handle routes
	handleConnections(&app, matchesApi, wsManager)

	log.Println("Server is started at port", app.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(app.Port), nil))
}

func handleConnections(app *config.App, matchesList *matches.List, wsManager *ws.Manager) {
	// Serve static files from the "static" directory
	http.Handle("/", http.FileServer(http.Dir("./ui")))

	// TODO  WebSocket endpoint
	http.HandleFunc("/ws", wsManager.ServeWS)
	http.HandleFunc("/ws-debug", wsManager.DebugClients)

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
