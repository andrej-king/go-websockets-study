package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

// WebSocket upgrader configuration
var upgrader = websocket.Upgrader{
	// Allow all origins for testing. In production, ensure you check the origin.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Client Struct to represent a connected client
type Client struct {
	Conn         *websocket.Conn // WebSocket connection
	sync.RWMutex                 // Mutex to synchronize writes
}

// Message Define a message structure
type Message struct {
	Type string `json:"type"` // Type of event (e.g., "regular_update", "request_response")
	Data string `json:"data"` // Payload data
}

// Subscribers keep subscribers by for specific events
type Subscribers struct {
	Clients     map[*Client]bool // Track all connected clients
	ExtraUpdate map[*Client]bool // Clients subscribed to extra updates
}

type Broadcasts struct {
	Common      chan Message // Channel for broadcasting messages to all clients
	ExtraUpdate chan Message // Channel for broadcasting extra updates
}

// WebSocketConfig keep common config data (Subscribers, Broadcasts)
type WebSocketConfig struct {
	Subscribers *Subscribers
	Broadcasts  *Broadcasts
}

var WebSocketConfigData *WebSocketConfig

func New() {
	WebSocketConfigData = &WebSocketConfig{
		Subscribers: &Subscribers{Clients: make(map[*Client]bool), ExtraUpdate: make(map[*Client]bool)},
		Broadcasts:  &Broadcasts{Common: make(chan Message), ExtraUpdate: make(chan Message)},
	}
}

// SafeWriteJSON Helper function to send a message to a client
func (c *Client) SafeWriteJSON(msg Message) error {
	c.Lock()         // Lock the mutex before writing
	defer c.Unlock() // Unlock after writing

	return c.Conn.WriteJSON(msg)
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

	client := &Client{Conn: connection} // Create a new client

	// close and delete client from each connection map
	defer func() {
		connection.Close()
		delete(WebSocketConfigData.Subscribers.Clients, client)
		delete(WebSocketConfigData.Subscribers.ExtraUpdate, client)
	}()

	WebSocketConfigData.Subscribers.Clients[client] = true // Add client to the common map

	msg := Message{Type: "response", Data: "Welcome message"}
	client.SafeWriteJSON(msg)

	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			log.Println("read error:", err)

			break // Exit the loop if the client tries to close the connection or the connection with the interrupted client
		}

		// local debug message
		log.Println(string(message))
	}
}
