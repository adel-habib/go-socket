package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn    *websocket.Conn
	manager *Manager
	egress  chan []byte
}

type ClientList map[*Client]bool

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{conn: conn, manager: manager, egress: make(chan []byte)}
}

func (c *Client) ReadMessages() {
	defer c.conn.Close()
	for {
		log.Println("reading, this will block")
		msgType, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println(err)
			}
			break
		}
		log.Println(msgType)
		log.Println(string(msg))
		for wsClient := range c.manager.Clients {
			wsClient.egress <- msg
		}
	}
}
func (c *Client) WriteMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()
	for {
		select {
		case msg, ok := <-c.egress:
		if !ok {
			if err := c.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
				log.Println(err)
			}
			return
		}
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println(err)
		}
		log.Println("message sent")

		}
	}
}
