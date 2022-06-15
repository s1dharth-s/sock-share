package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }}

var clients = make(map[*websocket.Conn]bool)
var messages = make(chan []byte)
var message = make(map[string]interface{})

func wsHandler(w http.ResponseWriter, r *http.Request) {

	// generate client id
	// clientId := 1000
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Client Connected!")
	clients[ws] = true
	greeting := `<div id="idMessage" hx-swap-oob="beforeend" hx-preserve>` + "Hello from server!" + `<br></div>`
	if err = ws.WriteMessage(1, []byte(greeting)); err != nil {
		log.Println(err)
	}
	for {
		// msg := Message{}
		// err := ws.ReadJSON(msg)
		// if err != nil {
		// 	log.Println(err)
		// 	return
		_, msg, err := ws.ReadMessage()
		log.Println(string(msg))
		if err != nil {
			log.Println(err)
			return
		}
		messages <- msg
	}
}

func handleMessages() {
	for {
		msg := <-messages

		// this will send message to all the client
		for client := range clients {
			// client.WriteJSON(msg)
			if err := json.Unmarshal(msg, &message); err != nil {
				log.Println(err)
			}

			res := `<div id="idMessage" hx-swap-oob="beforeend" hx-preserve>` + message["chat_message"].(string) + `<br></div>`

			if err := client.WriteMessage(1, []byte(res)); err != nil {
				log.Print(err)
				return

			}
		}
	}
}

// func handleClients(ws *websocket.Conn) {
// 	for {
// 		// msg := Message{}
// 		// err := ws.ReadJSON(msg)
// 		// if err != nil {
// 		// 	log.Println(err)
// 		// 	return
// 		_, msg, err := ws.ReadMessage()
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}
// 		messages <- msg
// 		fmt.Println(msg)
// 	}
// }

// func reader(ws *websocket.Conn) {
// 	for {
// 		messageType, msg, err := ws.ReadMessage()
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}

// 		fmt.Println(string(msg))
// 		if err = json.Unmarshal(msg, &m); err != nil {
// 			log.Println(err)
// 		}

// 		res := `<div id="idMessage" hx-swap-oob="beforeend" hx-preserve>` + m["chat_message"].(string) + `<br></div>`

// 		if err := ws.WriteMessage(messageType, []byte(res)); err != nil {
// 			log.Print(err)
// 			return
// 		}
// 	}
// }

func main() {
	http.HandleFunc("/chatroom", wsHandler)
	http.Handle("/", http.FileServer(http.Dir("./")))
	go handleMessages()
	http.ListenAndServe("127.0.0.1:8000", nil)
}
