package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

// WebSocket upgrader configuration
var upgrader = websocket.Upgrader{
	// Allow all origins for testing. In production, ensure you check the origin.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Track connected clients
var clients = make(map[*websocket.Conn]bool)

// Channel for broadcasting messages to clients
var broadcast = make(chan string)

// Handle incoming WebSocket connections
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}
	defer ws.Close()

	// Add the client to the map
	clients[ws] = true

	// Listen for incoming messages from the client (optional)
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			// Remove the client if the connection is closed or an error occurs
			delete(clients, ws)
			break
		}
	}
}

// Broadcast messages to all connected clients
func handleMessages() {
	for {
		// Wait for a new message from the broadcast channel
		msg := <-broadcast
		// Send the message to all connected clients
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				// Remove clients with errors (e.g., disconnected clients)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	// WebSocket endpoint
	http.HandleFunc("/ws", handleConnections)

	// Start the message handling goroutine
	go handleMessages()

	// Send updated data every 5 seconds
	go func() {
		for {
			time.Sleep(2 * time.Second)
			// Example: send updated data to all clients
			broadcast <- "Updated data: " + time.Now().UTC().Format(time.RFC1123)
		}
	}()

	// Start the HTTP server
	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
