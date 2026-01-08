package main

import (
	"fmt"
	"net/http"
	"ws-chat/server/server"
)

func main() {
	var port = 8080
	server := server.NewServer()
	fmt.Println("ws-chat server")
	http.HandleFunc("/ws", server.HandleNewConnection)

	go server.DeliverMessages(server.Queue)

	fmt.Printf("Server listening on port %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Printf("Failed to bind address: %s\n", err)
	}
}
