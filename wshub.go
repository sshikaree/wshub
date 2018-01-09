package wshub

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// var (
//	wsHub *WSHub = NewWSHub()
// )

// WsHub is used to collect all active websocket
// connections to send broadcast messages
type WSHub struct {
	mu sync.Mutex
	// Registered clients
	connections map[*websocket.Conn]bool

	// Incoming messages to be sent broadcast
	// Broadcast chan []byte

	// Register connection
	// Register chan *wswrapper.WrappedConn

	// Unregister connection
	// Unregister chan *wswrapper.WrappedConn
}

// Send message to all registered websocket connections
func (hub *WSHub) SendBroadcast(messageType int, msg []byte) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	for conn := range hub.connections {
		// log.Println("Sending to ws connection")
		// _, err := conn.Write(msg)
		err := conn.WriteMessage(messageType, msg)
		// log.Println("DONE")
		if err != nil {
			log.Println(err)
			delete(hub.connections, conn)
			conn.Close()
		}
	}
}

// Add websocket connection to hub
func (hub *WSHub) Register(conn *websocket.Conn) {
	hub.mu.Lock()
	hub.connections[conn] = true
	hub.mu.Unlock()
}

// Remove websocket connection from hub
func (hub *WSHub) Unregister(conn *websocket.Conn) {
	hub.mu.Lock()
	delete(hub.connections, conn)
	hub.mu.Unlock()
}

// Returns nubmber of active connections
func (hub *WSHub) Len() int {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	return len(hub.connections)
}

func NewWSHub() *WSHub {
	hub := new(WSHub)
	hub.connections = make(map[*websocket.Conn]bool)
	// wh.Broadcast = make(chan []byte)
	// wh.Register = make(chan *wswrapper.WrappedConn)
	// wh.Unregister = make(chan *wswrapper.WrappedConn)

	return hub
}
