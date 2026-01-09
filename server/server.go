package server

import (
	"fmt"
	"net/http"

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

	client := Client{ID: request.RemoteAddr, WS: ws}
	server.Clients = append(server.Clients, client)
	server.Queue <- MessageEnvelope{Sender: "Server", Message: fmt.Sprintf("%s has joined the server.", client.ID)}

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Printf("ReadMessage failed: %s\n", err)

			var deadClient_idx int
			for client_idx, client := range server.Clients {
				if client.WS == ws {
					deadClient_idx = client_idx
					break
				}
			}

			deadClient := server.Clients[deadClient_idx]
			server.Clients = append(server.Clients[:deadClient_idx], server.Clients[deadClient_idx+1:]...)
			server.Queue <- MessageEnvelope{Sender: "Server", Message: fmt.Sprintf("%s has left the server.", deadClient.ID)}
			break
		}
		fmt.Printf("%s sent: %s\n", request.RemoteAddr, message)
		server.Queue <- MessageEnvelope{Sender: request.RemoteAddr, Message: string(message)}
	}
}

func (server *Server) DeliverMessages(queue <-chan MessageEnvelope) {
	var sentMessage string
	for envelope := range queue {
		for _, client := range server.Clients {
			// do not send the message back to the sender
			if client.ID == envelope.Sender {
				continue
			}

			if envelope.Sender == "Server" {
				sentMessage = envelope.Message
			} else {
				sentMessage = fmt.Sprintf("%s wrote: %s\n", client.ID, envelope.Message)
			}

			client.WS.WriteMessage(websocket.TextMessage, []byte(sentMessage))
		}
	}
}
