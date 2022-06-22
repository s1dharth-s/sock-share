package main

import "log"

type Hub struct {
	rooms      map[*Room]bool
	message    chan *Message
	register   chan *Room
	unregister chan *Room
}

func (h *Hub) Run() {
	log.Println("Starting Hub...")
	for {
		select {
		case room := <-h.register:
			h.rooms[room] = true
			log.Println("Room registered: ", room.roomID)

		case room := <-h.unregister:
			delete(h.rooms, room)
			close(room.msgChannel)
			log.Println("Room unregistered:", room.roomID)

		case msg := <-h.message:
			log.Println("Hub recieved message")
			for room := range h.rooms {
				if room.roomID == msg.rid {
					room.msgChannel <- msg
					log.Println("Sent message to client")
				}
			}
		}
	}
}

func newHub() *Hub {
	return &Hub{
		rooms:      make(map[*Room]bool),
		message:    make(chan *Message),
		register:   make(chan *Room),
		unregister: make(chan *Room),
	}
}
