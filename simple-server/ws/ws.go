package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Accepting all connections
}

type Server struct {
	clients map[*websocket.Conn]bool
}

func StartServer(port uint16, route string) *Server {
	server := Server{make(map[*websocket.Conn]bool)}

	go func() {
		http.HandleFunc(route, server.handleConnections)

		err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

		if err != nil {
			fmt.Println("ListenAndServe:", err)
		}
	}()

	return &server
}

func (server *Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a WebSocket
	connection, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("upgrade:", err)
		return
	}

	// close and delete connection
	defer func() {
		connection.Close()
		delete(server.clients, connection)
	}()

	server.clients[connection] = true // Save the connection using it as a key

	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			fmt.Println("read error:", err)

			break // Exit the loop if the client tries to close the connection or the connection with the interrupted client
		}

		go server.sendMessageToAllClient(message) //  Send messages to all clients

		// local debug message
		fmt.Println(server.clients, string(message))
	}
}

// func (server *Server) sendMessageToAllClient(sender *websocket.Conn, message []byte) {
func (server *Server) sendMessageToAllClient(message []byte) {
	for client := range server.clients {
		//if client == sender {
		//	continue
		//}

		if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
			fmt.Println("broadcast error:", err)
			client.Close()
			delete(server.clients, client)
		}
	}
}
