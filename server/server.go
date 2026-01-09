package server

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/gorilla/websocket"
)

type MessageEnvelope struct {
	Sender  string
	Message string
}

type Server struct {
	Clients []Client
	Queue   chan MessageEnvelope
}

func NewServer() Server {
	return Server{Queue: make(chan MessageEnvelope)}
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

	server.Clients = append(server.Clients, Client{ID: request.RemoteAddr, WS: ws})

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Printf("ReadMessage failed: %s\n", err)
			server.Clients = slices.DeleteFunc(server.Clients, func(client Client) bool {
				return client.WS == ws
			})
			break
		}
		fmt.Printf("%s sent: %s\n", request.RemoteAddr, message)
		server.Queue <- MessageEnvelope{Sender: request.RemoteAddr, Message: string(message)}
	}
}

func (server *Server) DeliverMessages(queue <-chan MessageEnvelope) {
	for envelope := range queue {
		for _, client := range server.Clients {
			// do not send the message back to the sender
			if client.ID == envelope.Sender {
				continue
			}
			sentMessage := fmt.Sprintf("%s wrote: %s\n", client.ID, envelope.Message)
			client.WS.WriteMessage(websocket.TextMessage, []byte(sentMessage))
		}
	}
}
