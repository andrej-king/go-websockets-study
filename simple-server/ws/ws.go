package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Accepting all connections
}

type Server struct {
	Clients map[*websocket.Conn]bool
}

func Start() *Server {
	server := Server{make(map[*websocket.Conn]bool)}
	http.HandleFunc("/ws", server.handler)

	return &server
}

func (server *Server) handler(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a WebSocket
	connection, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	// close and delete connection
	defer func() {
		connection.Close()
		delete(server.Clients, connection)
	}()

	server.Clients[connection] = true // Save the connection using it as a key

	fmt.Println("Client Connected")
	err = connection.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println("Hello msg:", err)
		return
	}

	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			log.Println("read error:", err)

			break // Exit the loop if the client tries to close the connection or the connection with the interrupted client
		}

		go server.sendMessageToAllClient(message) //  Send messages to all clients

		// local debug message
		fmt.Println(server.Clients, string(message))
	}
}

// func (server *Server) sendMessageToAllClient(sender *websocket.Conn, message []byte) {
func (server *Server) sendMessageToAllClient(message []byte) {
	for client := range server.Clients {
		//if client == sender {
		//	continue
		//}

		if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("broadcast error:", err)
			client.Close()
			delete(server.Clients, client)
		}
	}
}
