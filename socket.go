package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Manager struct {
	Clients ClientList
	sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{Clients: make(ClientList)}
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection")
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := NewClient(conn, m)

	m.AddClient(client)

	go client.ReadMessages()

}
func (m *Manager) AddClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	m.Clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.Clients[client]; ok {
		client.conn.Close()
		delete(m.Clients, client)
	}
}
