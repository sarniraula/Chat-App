package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

//We'll need to upgrade the connction
//so define an upgrader which will require
//read and write buffer size ie Step 2

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	//we'll need to check the origin of our connection
	//this will allow us to make requests from our React
	//development server to here.
	//For now, we'll do no checking and just allow any connection ie Step 1
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//define a reader which will listen for
//new messages being sent to our WebSocket
//endpoint ie step 3

func reader(conn *websocket.Conn) {
	messageType, msg, err := conn.ReadMessage()

	if err != nil {
		log.Println(err)
		return
	}

	//print out that message for clarity
	fmt.Println(string(msg))

	if err := conn.WriteMessage(messageType, msg); err != nil {
		log.Println(err)
		return
	}
}

//define our websocket endpoint
func serveWS(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host)

	//upgrade this connection to a websocket connection (Step 2)
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	//listen indefinitely for new messages coming
	//through our websocket connection
	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "simple server")
	})

	//make our '/ws' endpoint to the 'serveWS' function
	http.HandleFunc("/ws", serveWS)
}

func main() {
	fmt.Println("Chat App v0.01")
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}
