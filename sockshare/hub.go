package sockshare

import "log"

type Hub struct {
	register   chan *Client
	unregister chan *Client
	message    chan *ChatMessage
	clients    map[*Client]bool
}

func (h *Hub) Run() {
	log.Println("Started Hub Service...")
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Println("A client has registered!")

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			log.Println("A client has unregistered!")

		case msg := <-h.message:
			log.Println("Hub received message from client id: ", msg.id)
			for client := range h.clients {
				client.send <- msg
			}
			log.Println("Sent message to client!")
		}
	}
}

func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		message:    make(chan *ChatMessage),
		clients:    make(map[*Client]bool),
	}
}
