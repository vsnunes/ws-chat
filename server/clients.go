package server

import "github.com/gorilla/websocket"

// Client struct to store info about connected clients
type Client struct {
	ID string
	WS *websocket.Conn
}
