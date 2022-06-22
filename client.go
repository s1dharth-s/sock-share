package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Room struct {
	count      int
	roomID     string
	cons       map[int]*websocket.Conn
	msgChannel chan *Message
}

type Message struct {
	Chat string `json:"chat_message"`
	rid  string
	cid  string
}

var Rooms = make(map[string]*Room)

func createRoomID(hub *Hub, c *gin.Context) {
	roomNo, ok := c.GetQuery("roomcode")

	if !ok {
		// TO-DO: Generate Random Roomcode
		num := 1000 + rand.Intn(8999)
		roomNo = strconv.Itoa(num)
		Rooms[roomNo] = &Room{
			count:      1,
			roomID:     roomNo,
			cons:       make(map[int]*websocket.Conn),
			msgChannel: make(chan *Message),
		}

		room := Rooms[roomNo]
		hub.register <- room
		go room.handleReads(hub)
		go room.handleWrites(hub)

		log.Println("New room created with roomID: ", roomNo)
	} else {
		log.Println("Room ID is: ", roomNo)
		Rooms[roomNo].count += 1
	}

	c.HTML(http.StatusOK, "chat.html", gin.H{"RoomNo": roomNo})
}

func connectRoom(c *gin.Context) {
	roomNo, ok := c.GetQuery("roomno")
	if ok {
		log.Println("Making new connection for ", roomNo)
	}
	r := Rooms[roomNo]
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error1", err)
	}
	r.cons[r.count] = ws
	log.Println("Client connected: ", r.count)
	// go handleReads(r, hub)
	// go handleWrites(r, hub)
}

func (r *Room) handleReads(hub *Hub) {
	defer func() {
		hub.unregister <- r
		for _, con := range r.cons {
			con.Close()
		}
	}()

	msg := &Message{}

	for {
		for cid, conn := range r.cons {
			err := conn.ReadJSON(msg)
			if err != nil {
				log.Println(err)
				break
			}
			msg.cid = strconv.Itoa(cid)
			msg.rid = r.roomID
			log.Println("Message: ", msg.rid, " ", msg.Chat, " ", msg.cid)
			hub.message <- msg
		}
	}
}

func (r *Room) handleWrites(hub *Hub) {
	for {
		msg := <-r.msgChannel

		if msg == nil {
			continue
		}

		sendMsg := `<div hx-swap-oob="beforeend:#content">` + msg.Chat + `<br></div>`

		for _, conn := range r.cons {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(sendMsg)); err != nil {
				log.Print("ERROR: ", err)
				return
			}
		}
	}
}
