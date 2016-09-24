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
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
