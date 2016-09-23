package main

import "github.com/gorilla/websocket"

// client is a user of chatting
type client struct {
	// socket is a websocket for this client
	socket *websocket.Conn
	// send is a channel for message
	send chan []byte
	// room is a chat room for this client
	room *room
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
