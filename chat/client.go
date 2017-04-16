package main

import (
	"time"

	"github.com/gorilla/websocket"
)

// client is a user of chatting
type client struct {
	// socket is a websocket for this client
	socket *websocket.Conn
	// send is a channel for message
	send chan *message
	// room is a chat room for this client
	room *room
	// userData have user information
	userData map[string]interface{}
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			c.room.forward <- msg
		} else {
			return
		}
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			return
		}
	}
}
