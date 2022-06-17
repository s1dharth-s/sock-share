package sockshare

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	hub  *Hub
	send chan *ChatMessage
	chat *ChatMessage
	id   int
}

type ChatMessage struct {
	Message string `json:"chat_message"`
	id      int
}

var id = 1

// var chats = &ChatMessage{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }}

func (c *Client) readMessages() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		err := c.conn.ReadJSON(c.chat)
		if err != nil {
			log.Println(err)
			break
		}
		log.Println("Message: ", c.chat.Message)
		c.hub.message <- c.chat
	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		msg := <-c.send
		//log.Println("Writing message: ", string(msg.Message))

		sendMsg := `<div hx-swap-oob="beforeend:#content">` + msg.Message + `<br></div>`
		if err := c.conn.WriteMessage(websocket.TextMessage, []byte(sendMsg)); err != nil {
			log.Print(err)
			return
		}
	}
}

func HandleSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error1", err)
	}
	log.Println("Client connected")
	client := &Client{conn: ws,
		hub:  hub,
		send: make(chan *ChatMessage),
		chat: &ChatMessage{id: id},
		id:   id}

	id += 1
	client.hub.register <- client

	go client.readMessages()
	go client.writeMessages()
}
