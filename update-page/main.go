package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go_websocket/update-page/api"
	"log"
	"net/http"
	"strconv"
	"time"
)

var port = 8000
var updateMatchesInterval = 7 * time.Second

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/api").Handler(makeApiHandler())
	//r.PathPrefix("/ws").HandlerFunc()
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend")))
	http.Handle("/", r)

	log.Printf("Going to listen on port %d\n", port)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(port), nil))
}

// makeApiHandler api routes
func makeApiHandler() http.Handler {
	// init match updates
	go api.Init(updateMatchesInterval)

	r := mux.NewRouter()

	// api home page
	r.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("{}")
	}).Methods("GET")

	// matches
	r.HandleFunc("/api/matches", func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(api.MatchList) // or use json.Marshal
	}).Methods("GET")

	return r
}
