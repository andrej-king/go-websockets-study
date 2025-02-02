package main

import (
	"log"
	"net/http"
	"strconv"
)

var port = 8080

func main() {
	Start()

	log.Printf("Going to listen on port %d\n", port)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(port), nil))
}
