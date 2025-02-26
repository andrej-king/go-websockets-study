package ws

import (
	"errors"
	"github.com/gorilla/websocket"
	"go_websocket/update-page/config"
	"net/http"
	"sync"
)

var ErrEventNotSupported = errors.New("this event type is not supported")

// WebSocket upgrader configuration
var websocketUpgrader = websocket.Upgrader{
	CheckOrigin: checkOrigin, // apply the Origin Checker
	//ReadBufferSize:  1024,
	//WriteBufferSize: 1024,
}

func checkOrigin(r *http.Request) bool {
	// Grab the request origin
	//origin := r.Header.Get("Origin")

	return true

	// TODO check host and need
	//switch origin {
	//case "http://localhost:8000":
	//	return true
	//default:
	//	return false
	//}
}

type Manager struct {
	// clients list for active clients
	clients *Subscribers

	// Using a syncMutex here to be able to lock state before editing clients
	// Could also use Channels to block
	sync.RWMutex

	// handlers are functions that are used to handle Events
	handlers map[string]EventHandler

	// LiveOddsChan channel keep messages for live odds subscribers
	LiveOddsChan chan Event
}

func NewManager(app *config.App) *Manager {
	m := &Manager{
		clients: &Subscribers{
			all:      make(map[*Client]bool),
			liveOdds: make(map[*Client]bool),
		},
		handlers:     make(map[string]EventHandler),
		LiveOddsChan: make(chan Event),
	}

	m.setupEventHandlers()

	return m
}

// setupEventHandlers add all handlers
func (m *Manager) setupEventHandlers() {
	m.handlers[EventSubscribe] = SubscriptionEventHandler
}

// routeEvent is used to make sure the correct event goes into the correct handler
func (m *Manager) routeEvent(event Event, client *Client) error {
	// Check if Handler is present in Map
	if handler, ok := m.handlers[event.Type]; ok {
		// Execute the handler and return any err
		if err := handler(event, client); err != nil {
			return err
		}

		return nil
	}

	return ErrEventNotSupported
}

// addClient add client to specific subscribe list
func (m *Manager) addClient(client *Client, clientType string) {
	m.Lock()
	defer m.Unlock()

	m.clients.all[client] = true

	switch clientType {
	case EventLiveOdds:
		m.clients.liveOdds[client] = true
	}
}

// removeClient remove client from every subscribe list
func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	// close client connection anyway
	client.connection.Close()

	// remove client from all connections list
	if _, ok := m.clients.all[client]; ok {
		delete(m.clients.all, client)
	}

	// remove client from live odd subscribers list
	if _, ok := m.clients.liveOdds[client]; ok {
		delete(m.clients.liveOdds, client)
	}
}

// SubscribersDataHandler send message to each player in channel
func (m *Manager) SubscribersDataHandler(event Event) {
	switch event.Type {
	case EventLiveOdds:
		for client := range m.clients.liveOdds {
			client.egress <- event
		}
	}
}
