package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

// WebSocket upgrader configuration
var upgrader = websocket.Upgrader{
	// Allow all origins for testing. In production, ensure you check the origin.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Struct to represent a connected client
type Client struct {
	Conn         *websocket.Conn // WebSocket connection
	sync.RWMutex                 // Mutex to synchronize writes
}

var clients = make(map[*Client]bool)                // Track connected clients
var extraUpdateSubscribers = make(map[*Client]bool) // Clients subscribed to extra updates
var broadcast = make(chan Message)                  // Channel for broadcasting messages to clients
var extraUpdates = make(chan Message)               // Channel for broadcasting extra updates

// Define a message structure
type Message struct {
	Type string `json:"type"` // Type of event (e.g., "regular_update", "request_response")
	Data string `json:"data"` // Payload data
}

// Helper function to send a message to a client
func (c *Client) SafeWriteJSON(msg Message) error {
	c.Lock()         // Lock the mutex before writing
	defer c.Unlock() // Unlock after writing

	return c.Conn.WriteJSON(msg)
}

// Handle incoming WebSocket connections
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}
	defer ws.Close()

	client := &Client{Conn: ws} // Create a new client
	clients[client] = true      // Add client to the map

	// Listen for incoming messages from the client
	for {
		var msg Message
		err := ws.ReadJSON(&msg)

		if err != nil {
			fmt.Println("Error reading message:", err)

			// Remove the client if the connection is closed or an error occurs
			delete(clients, client)

			// Remove from extra updates if client disconnects
			delete(extraUpdateSubscribers, client)

			break
		}

		// Handle different types of events
		switch msg.Type {
		case "request_response":
			// Respond to a specific client request
			response := Message{
				Type: "response",
				Data: "Response to your request: " + msg.Data,
			}

			client.SafeWriteJSON(response) // Respond only to the requesting client
		case "start_extra_updates":
			// Subscribe the client to extra updates
			extraUpdateSubscribers[client] = true

			fmt.Println("Client subscribed to extra updates")
		case "stop_extra_updates":
			// Remove extra updates subscriber
			delete(extraUpdateSubscribers, client)

			fmt.Println("Client unsubscribed from extra updates")
		}
	}
}

// Broadcast messages to all connected clients
func handleMessages() {
	for {
		msg := <-broadcast // Wait for a new message from the broadcast channel

		for client := range clients {
			err := client.SafeWriteJSON(msg) // Send the message as JSON

			if err != nil {
				fmt.Println("Error sending message:", err)

				// Remove clients with errors (e.g., disconnected clients)
				client.Conn.Close()

				delete(clients, client)
				delete(extraUpdateSubscribers, client) // Remove from extra updates if disconnected
			}
		}
	}
}

func handleExtraUpdates() {
	for {
		msg := <-extraUpdates // Wait for extra update messages channel

		for client := range extraUpdateSubscribers {
			err := client.SafeWriteJSON(msg) // Send the extra update message

			if err != nil {
				fmt.Println("Error sending extra update:", err)

				client.Conn.Close()
				delete(extraUpdateSubscribers, client) // Remove from extra updates on error
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

	// Start the message handling goroutines
	go handleMessages()
	go handleExtraUpdates()

	// Send regular updates to all clients
	go func() {
		for {
			time.Sleep(3 * time.Second) // Custom interval for regular updates

			// Broadcast regular updates to all clients
			broadcast <- Message{
				Type: "regular_update",
				Data: "Current time: " + time.Now().UTC().Format(time.RFC1123),
			}
		}
	}()

	// Send extra updates to subscribed clients with a different interval
	go func() {
		for {
			time.Sleep(2 * time.Second) // Custom interval for extra updates

			extraUpdates <- Message{
				Type: "extra_update",
				Data: "Extra update: " + time.Now().UTC().Format(time.RFC1123),
			}
		}
	}()

	// Start the HTTP server
	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
