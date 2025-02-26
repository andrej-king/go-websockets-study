package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

type Subscribers struct {
	all      map[*Client]bool // all connections
	liveOdds map[*Client]bool // live odd subscribers
}

// Client is a websocket client, basically a frontend visitor
type Client struct {
	connection *websocket.Conn // the websocket connection
	manager    *Manager        // is the manager used to manage the client
	egress     chan Event      // used to avoid concurrent writes on the WebSocket
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
	}
}

// readMessages will start the client to read messages and handle them appropriately
// This is supposed to be run as a goroutine
func (c *Client) readMessages() {
	defer func() {
		// graceful Close the Connection once this function is done
		c.manager.removeClient(c)
	}()

	// Set Max Size of Messages in Bytes
	// TODO check needed
	//c.connection.SetReadLimit(512)

	//	// welcome message for every new client
	//	msg := Event{Type: "response", Payload: "Welcome message"}
	//	c.writeMessages(msg)
	//var test Event
	//WelcomeMessageEventHandler(test, c)

	// Loop Forever
	for {
		_, payload, err := c.connection.ReadMessage()

		if err != nil {
			// If connection is closed, we will receive an error here
			// Log only strange errors
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}

			break // Break the loop to close conn & Cleanup
		}

		// Marshal incoming data into an Event struct
		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error marshalling message: %v", err)
			break // Breaking the connection here might be harsh
		}

		// Route event
		if err := c.manager.routeEvent(request, c); err != nil {
			log.Println("Error handling Message: ", err)
		}
	}
}

// writeMessages is a process that listens for new messages to output to the Client
func (c *Client) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	// Loop Forever
	for {
		select {
		case message, ok := <-c.egress:
			// Ok will be false Incase the egress channel is closed
			if !ok {
				// Manager has closed this connection channel, so communicate that to frontend
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					// Log that the connection is closed and the reason
					log.Println("connection closed: ", err)
				}

				// Return to close the goroutine
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return // closes the connection
			}

			// Write a Regular text message to the connection
			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println(err)
			}
		}
	}
}
