package ws

import (
	"fmt"
	"log"
	"net/http"
)

//const (
//ClientTypeCommon   = "Common"
//ClientTypeLiveOdds = "LiveOdds"
//)

// WebSocket upgrader configuration
//var upgrader = websocket.Upgrader{
//	// Allow all origins for testing. In production, ensure you check the origin.
//	CheckOrigin: func(r *http.Request) bool { return true },
//}

// Client Struct to represent a connected client
//type Client struct {
//	connection *websocket.Conn // WebSocket connection
//
//}

// Message Define a message structure
//type Message struct {
//	Type string `json:"type"` // Type of event (e.g., "regular_update", "request_response")
//	Data string `json:"data"` // Payload data
//}

// Subscribers keep subscribers by for specific events
//type Subscribers struct {
//	sync.RWMutex // Mutex to synchronize writes
//
//	Clients  map[*Client]bool // Track all connected clients
//	LiveOdds map[*Client]bool // Clients subscribed to extra updates
//}

//type Broadcasts struct {
//	sync.RWMutex // Mutex to synchronize writes
//
//	Common   chan Message // Channel for broadcasting messages to all clients
//	LiveOdds chan Message // Channel for broadcasting extra updates
//}

// WebSocketConfig keep common config data (Subscribers, Broadcasts)
//type WebSocketConfig struct {
//Subscribers *Subscribers
//Broadcasts  *Broadcasts
//}

//var WebSocketConfigData *WebSocketConfig

//func New() {
//	WebSocketConfigData = &WebSocketConfig{
//		Subscribers: &Subscribers{
//			Clients:  make(map[*Client]bool),
//			LiveOdds: make(map[*Client]bool),
//		},
//		Broadcasts: &Broadcasts{
//			Common:   make(chan Message),
//			LiveOdds: make(chan Message),
//		},
//	}
//}

// SafeWriteJSON Helper function to send a message to a client
//func (c *Client) SafeWriteJSON(msg Message) error {
//	c.Lock()         // Lock the mutex before writing
//	defer c.Unlock() // Unlock after writing
//
//	return c.Conn.WriteJSON(msg)
//}

//func HandleConnections(w http.ResponseWriter, r *http.Request) {
//	// Upgrade the HTTP connection to a WebSocket connection
//	connection, err := upgrader.Upgrade(w, r, nil)
//	if err != nil {
//		log.Println("Error upgrading connection:", err)
//		return
//	}
//
//	client := &Client{Conn: connection} // Create a new client
//
//	// close and delete client from each connection map
//	defer func() {
//		connection.Close()
//		delete(WebSocketConfigData.Subscribers.Clients, client)
//		delete(WebSocketConfigData.Subscribers.LiveOdds, client)
//	}()
//
//	WebSocketConfigData.Subscribers.Clients[client] = true // Add client to the common map
//
//	// welcome message for every new client
//	msg := Message{Type: "response", Data: "Welcome message"}
//	client.SafeWriteJSON(msg)
//
//	for {
//		mt, message, err := connection.ReadMessage()
//
//		if err != nil || mt == websocket.CloseMessage {
//			// If connection is closed, we will receive an error here
//			// Log only strange errors
//			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
//				log.Printf("error reading message: %v", err)
//			}
//
//			break // Break the loop to close conn & Cleanup
//		}
//
//		// local debug message
//		log.Println(string(message))
//	}
//}

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
