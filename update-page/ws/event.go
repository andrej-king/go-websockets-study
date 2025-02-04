package ws

import (
	"encoding/json"
	"fmt"
	"time"
)

// Event is the Messages sent over the websocket
// Used to differ between different actions
type Event struct {
	// Type is the message type sent
	Type string `json:"type"`
	// Payload is the data Based on the Type
	Payload json.RawMessage `json:"payload"`
}

// EventHandler is a function signature that is used to affect messages on the socket and triggered
// depending on the type
type EventHandler func(event Event, client *Client) error

const (
	// EventSendOdds is event when send odds
	EventSendOdds = "send_odds"
)

type NewMessageEvent struct {
	Message string    `json:"message"`
	Sent    time.Time `json:"sent"`
}

func SendOddsHandler(event Event, c *Client) error {
	var broadMessage NewMessageEvent
	broadMessage.Sent = time.Now().UTC()
	broadMessage.Message = "test message from server" // TODO replace to matches from API

	data, err := json.Marshal(broadMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventSendOdds

	// Broadcast to all other Clients
	for client := range c.manager.clients {
		// Only send to clients inside the same chatroom
		//if client.room == c.room {
		client.egress <- outgoingEvent
		//}
	}

	return nil
}
