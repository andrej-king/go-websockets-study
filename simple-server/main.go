package main

// https://dev.to/davidnadejdin/simple-server-on-gorilla-websocket-52h7
// https://dev.to/neelp03/using-websockets-in-go-for-real-time-communication-4b3l
// https://github.com/gorilla/websocket/blob/main/examples/echo/
// https://eli.thegreenplace.net/2019/on-concurrency-in-go-http-servers/

import (
	"go_websocket/simple-server/ws"
	"log"
	"net/http"
	"strconv"
)

var port = 8080

func main() {
	ws.Start()

	log.Printf("Going to listen on port %d\n", port)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(port), nil))
}
