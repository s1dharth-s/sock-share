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
	msg        *Message
}

type Message struct {
	// mu   sync.Mutex
	Chat string `json:"chat_message"`
	rid  string
	cid  string
}

// var msg = &Message{}

var Rooms = make(map[string]*Room)

func createRoomID(c *gin.Context) {
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
			msg:        &Message{},
		}

		room := Rooms[roomNo]
		// hub.register <- room
		// go room.handleReads(hub)
		go room.handleWrites()

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
	} else {
		panic("Could not find room number!")
	}
	r := Rooms[roomNo]
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error1", err)
	}
	r.cons[r.count] = ws
	log.Println("Client connected: ", r.count)

	for {
		err := ws.ReadJSON(r.msg)
		if err != nil {
			log.Println(err)
			break
		}
		r.msg.cid = strconv.Itoa(r.count)
		r.msg.rid = r.roomID
		log.Println("Message: ", r.msg.rid, " ", r.msg.Chat, " ", r.msg.cid)
		// hub.message <- r.msg
		r.msgChannel <- r.msg
	}
}

func (r *Room) handleWrites() {
	for {
		msg := <-r.msgChannel
		if msg == nil {
			log.Println("nil message")
			continue
		}
		log.Println("Received message: ", msg.Chat)

		sendMsg := `<div hx-swap-oob="beforeend:#content">` + msg.Chat + `<br></div>`

		for _, conn := range r.cons {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(sendMsg)); err != nil {
				log.Print("ERROR: ", err)
				// return
			}
			log.Println("Sent message!")
		}
	}
}
