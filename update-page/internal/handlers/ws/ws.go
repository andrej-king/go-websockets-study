package ws

import (
	"fmt"
	"log"
	"net/http"
)

func (m *Manager) DebugClients(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "all: ", len(m.clients.all), ", live: ", len(m.clients.liveOdds))
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("Client connected")

	// Begin by upgrading the HTTP request
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	// Create New Client
	client := NewClient(conn, m)
	m.addClient(client, "")

	go client.readMessages()
	go client.writeMessages()
}
