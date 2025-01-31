package main

// https://dev.to/davidnadejdin/simple-server-on-gorilla-websocket-52h7

import (
	"fmt"
	"web_socket_server/ws"
)

func main() {
	ws.StartServer(8080, messageHandler)

	// wait input for continue work server
	var input string
	fmt.Scanln(&input)
}

func messageHandler(message []byte) {
	fmt.Println(string(message))
}
