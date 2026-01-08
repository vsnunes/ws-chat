package server

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/gorilla/websocket"
)

type Server struct {
	Clients []*websocket.Conn
	Queue   chan string
}

func NewServer() Server {
	return Server{Queue: make(chan string)}
}

func (server *Server) HandleNewConnection(writer http.ResponseWriter, request *http.Request) {
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

	server.Clients = append(server.Clients, ws)

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Printf("ReadMessage failed: %s\n", err)
			server.Clients = slices.DeleteFunc(server.Clients, func(client *websocket.Conn) bool {
				return client == ws
			})
			break
		}
		fmt.Printf("Received: %s\n", message)
		server.Queue <- string(message)
	}
}

func (server *Server) DeliverMessages(queue <-chan string) {
	for message := range queue {
		for _, client := range server.Clients {
			client.WriteMessage(websocket.TextMessage, []byte(message))
		}
	}
}
