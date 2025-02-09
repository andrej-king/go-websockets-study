package main

import (
	"encoding/json"
	"go_websocket/update-page/internal/handlers/api/matches"
	"log"
	"net/http"
	"strconv"
)

type App struct {
	Port          int                    // start app port
	LiveMatchList *map[int]matches.Match // fill by api request with interval
}

func main() {
	liveMatchList := make(map[int]matches.Match)
	app := App{
		Port:          8080,
		LiveMatchList: &liveMatchList,
	}

	app.handleConnections()

	log.Println("Server is started at port", app.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(app.Port), nil))
}

func (a *App) handleConnections() {
	// Serve static files from the "static" directory
	http.Handle("/", http.FileServer(http.Dir("./public")))

	// TODO  WebSocket endpoint
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(a)
	})

	// TODO live matches endpoint
	http.HandleFunc("/api/matches/live", func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(a.LiveMatchList)
	})
}
