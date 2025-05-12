package ws

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Message represents a message to be broadcasted
type Message struct {
	Room    string      `json:"room"`
	Payload interface{} `json:"payload"`
}

// Hub manages active WebSocket connections
type Hub struct {
	Clients   map[*websocket.Conn]string
	Lock      sync.Mutex
	Broadcast chan Message
}

var hub = Hub{
	Clients:   make(map[*websocket.Conn]string),
	Broadcast: make(chan Message),
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleConnections upgrades HTTP requests to WebSocket connections
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	hub.Lock.Lock()
	hub.Clients[conn] = ""
	hub.Lock.Unlock()

	defer func() {
		hub.Lock.Lock()
		delete(hub.Clients, conn)
		hub.Lock.Unlock()
		conn.Close()
	}()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		hub.Lock.Lock()
		hub.Clients[conn] = msg.Room
		hub.Lock.Unlock()
	}
}

// BroadcastMessages listens for messages and sends them to appropriate clients
func BroadcastMessages() {
	for {
		msg := <-hub.Broadcast
		hub.Lock.Lock()
		for client, room := range hub.Clients {
			if room == msg.Room {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Println("WebSocket write error:", err)
					client.Close()
					delete(hub.Clients, client)
				}
			}
		}
		hub.Lock.Unlock()
	}
}

// StartWebSocketServer initializes WebSocket handling
func StartWebSocketServer() {
	http.HandleFunc("/ws", HandleConnections)
	go BroadcastMessages()
	fmt.Println("WebSocket server started on /ws")
}

func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
}

// ===============================================================================

// package handlers

// import (
// 	"log"
// 	"net/http"
// 	"sync"

// 	"github.com/gin-gonic/gin"
// 	"github.com/gorilla/websocket"
// )

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool { return true },
// }

// type Message struct {
// 	Room    string      `json:"room"`
// 	Payload interface{} `json:"payload"`
// }

// type Client struct {
// 	conn  *websocket.Conn
// 	send  chan Message
// 	rooms map[string]bool
// }

// type Hub struct {
// 	clients    map[*Client]bool
// 	broadcast  chan Message
// 	register   chan *Client
// 	unregister chan *Client
// 	mu         sync.Mutex
// }

// var Hub = &Hub{
// 	broadcast:  make(chan Message),
// 	register:   make(chan *Client),
// 	unregister: make(chan *Client),
// 	clients:    make(map[*Client]bool),
// }

// func (h *Hub) Run() {
// 	for {
// 		select {
// 		case client := <-h.register:
// 			h.mu.Lock()
// 			h.clients[client] = true
// 			h.mu.Unlock()

// 		case client := <-h.unregister:
// 			h.mu.Lock()
// 			if _, ok := h.clients[client]; ok {
// 				close(client.send)
// 				delete(h.clients, client)
// 			}
// 			h.mu.Unlock()

// 		case msg := <-h.broadcast:
// 			h.mu.Lock()
// 			for client := range h.clients {
// 				if _, ok := client.rooms[msg.Room]; ok {
// 					select {
// 					case client.send <- msg:
// 					default:
// 						close(client.send)
// 						delete(h.clients, client)
// 					}
// 				}
// 			}
// 			h.mu.Unlock()
// 		}
// 	}
// }

// func HandleWebSocket(c *gin.Context) {
// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		log.Println("WebSocket upgrade error:", err)
// 		return
// 	}

// 	client := &Client{
// 		conn:  conn,
// 		send:  make(chan Message, 256),
// 		rooms: make(map[string]bool),
// 	}

// 	Hub.register <- client

// 	go client.writePump()
// 	go client.readPump()
// }

// func (c *Client) readPump() {
// 	defer func() {
// 		Hub.unregister <- c
// 		c.conn.Close()
// 	}()

// 	for {
// 		var msg struct {
// 			Action string `json:"action"`
// 			Room   string `json:"room"`
// 		}
// 		if err := c.conn.ReadJSON(&msg); err != nil {
// 			break
// 		}

// 		switch msg.Action {
// 		case "join":
// 			c.rooms[msg.Room] = true
// 		case "leave":
// 			delete(c.rooms, msg.Room)
// 		}
// 	}
// }

// func (c *Client) writePump() {
// 	defer c.conn.Close()
// 	for msg := range c.send {
// 		if err := c.conn.WriteJSON(msg); err != nil {
// 			break
// 		}
// 	}
// }
