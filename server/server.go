package server

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/gorilla/websocket"
)

type Server struct {
	Clients []Client
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

	server.Clients = append(server.Clients, Client{ID: "unknown", WS: ws})

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Printf("ReadMessage failed: %s\n", err)
			server.Clients = slices.DeleteFunc(server.Clients, func(client Client) bool {
				return client.WS == ws
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
			sentMessage := fmt.Sprintf("%s wrote: %s\n", client.ID, message)
			client.WS.WriteMessage(websocket.TextMessage, []byte(sentMessage))
		}
	}
}
