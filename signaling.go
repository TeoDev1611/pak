package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type SignalingServer struct {
	clients map[*websocket.Conn]bool
	mutex   sync.Mutex
}

func NewSignalingServer() *SignalingServer {
	return &SignalingServer{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (s *SignalingServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrade:", err)
		return
	}
	defer conn.Close()

	s.mutex.Lock()
	s.clients[conn] = true
	s.mutex.Unlock()

	fmt.Println("Nuevo cliente conectado al señalizador")

	for {
		var msg interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			s.mutex.Lock()
			delete(s.clients, conn)
			s.mutex.Unlock()
			break
		}

		// Reenviar el mensaje a todos los demás clientes (broadcast simple)
		s.mutex.Lock()
		for client := range s.clients {
			if client != conn {
				client.WriteJSON(msg)
			}
		}
		s.mutex.Unlock()
	}
}
