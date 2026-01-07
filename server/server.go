package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

func HandleNewConnection(writer http.ResponseWriter, request *http.Request) {
	// Send HTTP upgrade to upgrade HTTP to WebSockets
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	ws, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		fmt.Printf("Failed to upgrade HTTP connection to WS: %s\n", err)
		return
	}

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Printf("ReadMessage failed: %s\n", err)
			break
		}
		fmt.Printf("Received: %s\n", message)
	}
}

func main() {
	var port = 8080
	fmt.Println("ws-chat server")
	http.HandleFunc("/ws", HandleNewConnection)

	fmt.Printf("Server listening on port %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Printf("Failed to bind address: %s\n", err)
	}
}
