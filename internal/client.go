package internal

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func EstablishConnection(connName string, roomId string) (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8080", Path: "/room"}
	log.Printf("connection to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
		return nil, err
	}

	conn := Connection{
		Name:   connName,
		RoomId: roomId,
	}
	data, err := json.Marshal(conn)
	if err != nil {
		return nil, err
	}

	err = c.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		return nil, err
	}

	return c, nil
}
