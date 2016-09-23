package main

type room struct {
	// forward retain message
	forward chan []byte
	// join is a channel for client who wants to login
	join chan *client
	// leave is a channel for leave this room
	leave chan *client
	// clients is in this room
	clients map[*client]bool
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// join
			r.clients[client] = true
		case client := <-r.leave:
			// leave
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			// send all clients
			for client := range r.clients {
				select {
				case client.send <- msg:
				// send message
				default:
					// fail to send
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}
