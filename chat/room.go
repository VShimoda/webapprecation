package main

import (
	"log"
	"net/http"

	"github.com/VShimoda/webapprecation/trace"
	"github.com/gorilla/websocket"
)

type room struct {
	// forward retain message
	forward chan []byte
	// join is a channel for client who wants to login
	join chan *client
	// leave is a channel for leave this room
	leave chan *client
	// clients is in this room
	clients map[*client]bool
	// tracer is logging this chat room
	tracer trace.Tracer
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// join
			r.clients[client] = true
			r.tracer.Trace("new client joined")
		case client := <-r.leave:
			// leave
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("client left")
		case msg := <-r.forward:
			r.tracer.Trace(" -- sent to client")
			// send all clients
			for client := range r.clients {
				select {
				case client.send <- msg:
				// send message
				default:
					// fail to send
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace(" -- faled to send message. clean up clients")
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() {
		r.leave <- client
	}()
	go client.write()
	client.read()
}
