package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	RoomId  string `json:"room_id"`
	Content string `json:"content"`
}

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, buf, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		var message Message
		err = json.Unmarshal(buf, &message)
		if err != nil {
			log.Println("Error trying to parse message")
			err = c.WriteMessage(mt, []byte("Invalid message format"))
			if err != nil {
				log.Println("Error trying to send warning message")
				break
			}
			continue
		}

		log.Println("recv: ", message)
		err = c.WriteMessage(mt, []byte(message.Content))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

const port = "8080"

func main() {
	flag.Parse()
	http.HandleFunc("/echo", echo)
	log.Printf("Server running at port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("127.0.0.1:%s", port), nil))
}
