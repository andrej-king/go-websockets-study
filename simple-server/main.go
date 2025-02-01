package main

// https://dev.to/davidnadejdin/simple-server-on-gorilla-websocket-52h7
// https://dev.to/neelp03/using-websockets-in-go-for-real-time-communication-4b3l
// https://github.com/gorilla/websocket/blob/main/examples/echo/

import (
	"fmt"
	"go_websocket/simple-server/ws"
)

var port uint16 = 8080
var route string = "/ws"

func main() {
	fmt.Println(fmt.Sprintf("WebSocket server started on :%d%s", port, route))
	ws.StartServer(port, route)

	// wait input for continue work server
	var input string
	fmt.Scanln(&input)
}
