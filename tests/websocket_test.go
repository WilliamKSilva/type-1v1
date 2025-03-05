package client

import (
	"encoding/json"
	"log"
	"net/url"
	"testing"

	"github.com/WilliamKSilva/type-1v1/internal"
	"github.com/gorilla/websocket"
)

const mockedRoomId string = "21ca15d0-e346-4630-a240-773a828c31b3"

func TestWebsocketConnection(t *testing.T) {
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8080", Path: "/room"}
	log.Printf("connection to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	conn := internal.Connection{
		Name:   "Teste",
		RoomId: mockedRoomId,
	}
	data, err := json.Marshal(conn)
	if err != nil {
		t.Fatalf("Websocket connection: %s", err)
		return
	}

	c.WriteMessage(websocket.BinaryMessage, data)

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			t.Log("read: ", err)
			return
		}

		var connection internal.Connection
		err = json.Unmarshal(message, &connection)
		if err != nil {
			return
		}
		if connection.Name != "" {
			t.Log("Connection established:", connection.Name)
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			return
		}
	}
}
