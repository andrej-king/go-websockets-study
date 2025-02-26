package ws

import (
	"encoding/json"
	"log"
	"time"
)

// SubscriptionEvent subscription payload type
type SubscriptionEvent struct {
	Name string `json:"name"`
}

// NewMessageEvent common payload type
type NewMessageEvent struct {
	Message string    `json:"message"`
	Sent    time.Time `json:"sent"`
}

func SubscriptionEventHandler(event Event, client *Client) error {
	var subscription SubscriptionEvent
	if err := json.Unmarshal(event.Payload, &subscription); err != nil {
		log.Printf("error marshalling message: %v", err)
		return nil
	}

	// subscribe client
	client.manager.addClient(client, subscription.Name)

	return nil
}
