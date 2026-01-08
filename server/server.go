package server

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

		//reply back original message
		ws.WriteMessage(websocket.TextMessage, message)
	}
}
