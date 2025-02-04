package ws

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var ErrEventNotSupported = errors.New("this event type is not supported")

var websocketUpgrader = websocket.Upgrader{
	// Apply the Origin Checker
	CheckOrigin: checkOrigin,
	//ReadBufferSize:  1024,
	//WriteBufferSize: 1024,
}

func checkOrigin(r *http.Request) bool {
	// Grab the request origin
	origin := r.Header.Get("Origin")

	// TODO check host and need
	switch origin {
	case "http://localhost:8000":
		return true
	default:
		return false
	}
}

type Manager struct {
	// clients list for all active clients
	clients ClientList

	// Using a syncMutex here to be able to lock state before editing clients
	// Could also use Channels to block
	sync.RWMutex

	// handlers are functions that are used to handle Events
	handlers map[string]EventHandler
}

func NewManager() *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
	}

	m.setupEventHandlers()

	return m
}

func (m *Manager) DebugClients(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, len(m.clients))
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

	// Add the newly created client to the manager
	m.addClient(client)

	go client.readMessages()
	go client.writeMessages()
}

// setupEventHandlers configures and adds all handlers
func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendOdds] = SendOddsHandler
}

// routeEvent is used to make sure the correct event goes into the correct handler
func (m *Manager) routeEvent(event Event, c *Client) error {
	// Check if Handler is present in Map
	if handler, ok := m.handlers[event.Type]; ok {
		// Execute the handler and return any err
		if err := handler(event, c); err != nil {
			return err
		}

		return nil
	}

	return ErrEventNotSupported
}

func (m *Manager) addClient(client *Client) {
	// Lock before manipulate with map
	m.Lock()
	defer m.Unlock()

	// Add Client
	m.clients[client] = true
}
func (m *Manager) removeClient(client *Client) {
	// Lock before manipulate with map
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		// close connection
		client.connection.Close()

		// remove
		delete(m.clients, client)
	}
}
